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
	recurse := flag.BoolP("recursive", "r", false, "Return all files with given name")
	argroot := flag.String("root", ".", "The root directory of the search")
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

	if *help {
		flag.Usage()
		return
	}

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}

	target := args[0]

	if !*recurse {
		result, err := findTarget(target, *argroot, *dir, *cont, *abs)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	} else {
		results, err := findAllTargets(target, *argroot, *dir, *cont, *abs)
		if err != nil {
			log.Fatal(err)
		}

		for _, r := range results {
			fmt.Println(r)
		}
	}

}
