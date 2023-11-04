package main

import (
	flag "github.com/spf13/pflag"
)


func main() {
	f := getFlags()

	// If -h flag was passed or there is the wrong number of arguments
	if f.help || len(f.args) != 1 {
		flag.Usage()
		return
	}

	findAndPrintTarget(f)
}

