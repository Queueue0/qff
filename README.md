# Queueue's File Finder

This is just a basic file/directory finder I made.

## Installation
```
go install github.com/Queueue0/qff/cmd/qff@latest
```

## Usage
Output of `qff --help`:

```
Find the path to the specified file
Usage of qff:
	qff [flags] [filename]

Flags:
  -a, --absolute-path    Return absolute path
  -c, --containing-dir   Return path to the directory containing [filename]
  -d, --find-directory   Search for a directory instead of a file
  -h, --help             Displays usage information
  -r, --recursive        Return all files with given name
      --root string      The root directory of the search (default ".")
```
