package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

type flags struct {
	abs     bool
	cont    bool
	dir     bool
	recurse bool
	help    bool
	root    string
	args    []string
}

var progname = filepath.Base(os.Args[0])

func getFlags() flags {
	abs := flag.BoolP("absolute-path", "a", false, "Return absolute path")
	cont := flag.BoolP("containing-dir", "c", false, "Return path to the directory containing [filename]")
	dir := flag.BoolP("find-directory", "d", false, "Search for a directory instead of a file")
	recurse := flag.BoolP("recursive", "r", false, "Return all files with given name")
	root := flag.String("root", ".", "The root directory of the search")
	help := flag.BoolP("help", "h", false, "Displays usage information")

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

	f := flags{
		abs:     *abs,
		cont:    *cont,
		dir:     *dir,
		recurse: *recurse,
		help:    *help,
		root:    *root,
		args:    args,
	}

	return f
}
