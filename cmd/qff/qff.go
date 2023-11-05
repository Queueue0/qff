package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

func findAndPrintTarget(f flags) {
	target := f.args[0]

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
		results, err = makeRelative(root, results...)
		if err != nil {
			printErrorAndExit(err)
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
	go findConcurrently(target, root, dir, cont)

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

	return results[0], nil
}

func findConcurrently(target, root string, dir, cont bool) error {
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

			if d.IsDir() == dir && d.Name() == target {
				if !cont {
					s = path
				} else {
					s = filepath.Dir(path)
				}

				resultChan <- s
				quit <- true
				close(quit)
				return filepath.SkipDir
			}

			if d.IsDir() && path != root {
				wg.Add(1)
				go findConcurrently(target, path, dir, cont)
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

	var results []string
	go func() {
		for s := range resultChan {
			results = append(results, s)
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
		close(quit)
	}()

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

		if d.IsDir() == dir && d.Name() == target {
			if !cont {
				s = path
			} else {
				s = filepath.Dir(path)
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

func sanitizeRoot(root string) (string, error) {
	var err error
	if root == "." {
		root, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", err
	}

	return root, nil
}

func makeRelative(root string, args ...string) ([]string, error) {
	var err error
	for i, arg := range args {
		if len(arg) != 0 {
			args[i], err = filepath.Rel(root, arg)
			args[i] = "./" + args[i]
		}

		if err != nil {
			return nil, err
		}
	}

	return args, nil
}

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
}
