package main

import (
	"fmt"
	"os"
	"path/filepath"
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

func findTarget(target, root string, dir, cont bool) (string, error) {
	result := ""
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if (d.IsDir() == dir) && d.Name() == target {
			if !cont {
				result = path
			} else {
				result = filepath.Dir(path)
			}
			return filepath.SkipAll
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return result, nil
}

func findAllTargets(target, root string, dir, cont bool) ([]string, error) {
	var results []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if (d.IsDir() == dir) && d.Name() == target {
			if !cont {
				results = append(results, path)
			} else {
				results = append(results, filepath.Dir(path))
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
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
