package main

import (
	"os"
	"path/filepath"
)

func findTarget(target, root string, dir, cont, abs bool) (string, error) {
	root, err := sanitizeRoot(root)

	result := ""
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
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

	if !abs {
		result, err = filepath.Rel(root, result)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func findAllTargets(target, root string, dir, cont, abs bool) ([]string, error) {
	root, err := sanitizeRoot(root)

	var results []string
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
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

	if !abs {
		for i, result := range results {
			results[i], err = filepath.Rel(root, result)
		}
		if err != nil {
			return nil, err
		}
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
