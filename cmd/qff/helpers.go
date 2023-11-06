package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func printErrorAndExit(err error) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(1)
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

	root, err = filepath.Abs(root)
	if err != nil {
		return "", err
	}

	return root, nil
}

func makeRelative(root string, args ...string) ([]string, error) {
	var err error
	prefix := root + "/"
	if strings.HasSuffix(root, "/") {
		prefix = root
	}
	root, _ = sanitizeRoot(root)
	for i, arg := range args {
		if len(arg) != 0 {
			args[i], err = filepath.Rel(root, arg)
			args[i] = prefix + args[i]
		}

		if err != nil {
			return nil, err
		}
	}

	return args, nil
}

func parseWildCard(pattern string) string {
	var result strings.Builder
	for i, s := range strings.Split(pattern, "*") {
		if i > 0 {
			result.WriteString(".*")
		}

		result.WriteString(regexp.QuoteMeta(s))
	}
	return result.String()
}

func match(pattern, s string) bool {
	// We can safely ignore the error here because result will be false anyway
	result, _ := regexp.MatchString(parseWildCard(pattern)+"$", s)
	return result
}
