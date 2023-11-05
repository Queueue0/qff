package main

import (
	flag "github.com/spf13/pflag"
)


func main() {
	f := getFlags()

	if f.help || len(f.args) != 1 {
		flag.Usage()
		return
	}

	findAndPrintTarget(f)
}
