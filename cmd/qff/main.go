package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var progname = filepath.Base(os.Args[0])

func main() {
	abs := flag.BoolP("absolute-path", "a", false, "Return absolute path")
	cont := flag.BoolP("containing-dir", "c", false, "Return path to the directory containing [filename]")
	dir := flag.BoolP("find-directory", "d", false, "Search for a directory instead of a file")
	argroot := flag.StringP("root", "r", ".", "The root directory of the search")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Find the path to the specified file
Usage of %s:
	%s [flags] [filename]

Flags:
`, progname, progname)
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}

	target := args[0]

	root := *argroot
	var err error
	if root == "." {
		root, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		log.Fatal(err)
	}

	result := ""
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		if (d.IsDir() == *dir) && d.Name() == target {
			if !*cont {
				result = path
			} else {
				result = filepath.Dir(path)
			}
			return filepath.SkipAll
		}

		return nil
	})

	if result == "" {
		os.Exit(1)
	}

	if !*abs {
		result, err = filepath.Rel(root, result)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(result)
}
