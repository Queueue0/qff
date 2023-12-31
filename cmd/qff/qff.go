package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func findAndPrintTarget(f flags) {
	target := f.args[0]

	if strings.HasPrefix(target, ".") {
		prefix, err := os.Getwd()
		if err != nil {
			printErrorAndExit(err)
		}

		target = prefix + strings.TrimPrefix(target, ".")
	}

	if !strings.HasPrefix(target, "/") {
		target = "/" + target
	}

	if strings.HasSuffix(target, "/") {
		target = strings.TrimSuffix(target, "/")
	}

	root, err := sanitizeRoot(f.root)

	if err != nil {
		printErrorAndExit(err)
	}

	var results []string
	if !f.recurse {
		result, err := findTarget(target, root, f.dir, f.cont)
		if err != nil {
			printErrorAndExit(err)
		}

		results = append(results, result)
	} else {
		results, err = findAllTargets(target, root, f.dir, f.cont)
		if err != nil {
			printErrorAndExit(err)
		}
	}

	if !f.abs {
		results, err = makeRelative(f.root, results...)
		if err != nil {
			printErrorAndExit(err)
		}
	}

	if f.dir || f.cont {
		for i := range results {
			if results[i] != "" {
				results[i] += "/"
			}
		}
	}

	for _, r := range results {
		fmt.Println(r)
	}
}

var wg sync.WaitGroup
var resultChan = make(chan string, 1)
var errChan = make(chan error, 1)
var quit = make(chan bool, 1)

func findTarget(target, root string, dir, cont bool) (string, error) {
	wg.Add(1)
	go findConcurrently(target, root, root, dir, cont)

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return "", err
		}
	}

	var results []string
	for s := range resultChan {
		results = append(results, s)
	}
	sort.Strings(results)

	if len(results) > 0 {
		return results[0], nil
	}
	return "", nil
}

func findConcurrently(target, root, origRoot string, dir, cont bool) error {
	defer wg.Done()

	var s string
	visit := func(path string, d os.DirEntry, err error) error {
		select {
		case <-quit:
			return filepath.SkipDir
		default:
			if err != nil {
				return err
			}

			if d.IsDir() == dir && match(target, path) && path != origRoot {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}

				if !cont {
					s = absPath
				} else {
					s = filepath.Dir(absPath)
				}

				resultChan <- s
				quit <- true
				close(quit)
				return filepath.SkipDir
			}

			if d.IsDir() && path != root {
				wg.Add(1)
				go findConcurrently(target, path, origRoot, dir, cont)
				return filepath.SkipDir
			}

			return nil
		}
	}

	select {
	case <-quit:
		return nil
	default:
		err := filepath.WalkDir(root, visit)
		if err != nil {
			errChan <- err
			return err
		}
	}
	return nil
}

func findAllTargets(target, root string, dir, cont bool) ([]string, error) {
	wg.Add(1)
	go findAllConcurrently(target, root, dir, cont)

	var closed sync.WaitGroup
	resultSet := make(map[string]struct{})
	closed.Add(1)
	go func() {
		for s := range resultChan {
			resultSet[s] = struct{}{}
		}
		closed.Done()
	}()

	closed.Add(1)
	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
		close(quit)
		closed.Done()
	}()

	closed.Wait()

	results := make([]string, len(resultSet))

	i := 0
	for k := range resultSet {
		results[i] = k
		i++
	}

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	sort.Strings(results)

	return results, nil
}

func findAllConcurrently(target, root string, dir, cont bool) error {
	defer wg.Done()

	var s string
	visit := func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() == dir && match(target, path) {
			absPath, err := filepath.Abs(path)
			if err != nil {
				return err
			}

			if !cont {
				s = absPath
			} else {
				s = filepath.Dir(absPath)
			}

			resultChan <- s
		}

		if d.IsDir() && path != root {
			wg.Add(1)
			go findAllConcurrently(target, path, dir, cont)
			return filepath.SkipDir
		}

		return nil
	}

	err := filepath.WalkDir(root, visit)
	if err != nil {
		errChan <- err
		return err
	}

	return nil
}
