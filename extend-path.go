package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	flHelp        = flag.Bool("help", false, "print the usage message and exit")
	flDelete      = flag.Bool("delete", false, "delete entries from PATH")
	flAppend      = flag.Bool("append", false, "append entry to PATH")
	flPrepend     = flag.Bool("prepend", false, "prepend entry to PATH")
	flRelocate    = flag.Bool("relocate", false, "if entry is already present, relocate it")
	flValidate    = flag.Bool("validate", false, "skip non-existing directories")
	flEnvironment = flag.String("environment", "PATH", "environment variable to use as PATH")
)

func main() {
	flag.Parse()

	if false == run(flag.Args()) {
		os.Exit(1)
	}
}

func run(directories []string) bool {
	path := os.Getenv(*flEnvironment)
	paths := splitDirectoryPath(path)
	values := resolveDirectories(directories, *flValidate)

	if *flHelp {
		printUsage()
	} else if len(values) == 0 {
		printPath(paths)
	} else if *flDelete {
		printPath(deleteEntries(values, paths))
	} else if *flAppend {
		printPath(appendEntries(values, paths, *flRelocate))
	} else if *flPrepend {
		printPath(prependEntries(values, paths, *flRelocate))
	} else {
		return verifyEntries(values, paths)
	}

	return true
}

func printPath(directories []Directory) {
	var buffer bytes.Buffer
	var first = true

	for _, directory := range directories {
		if first {
			first = false
		} else {
			buffer.WriteRune(filepath.ListSeparator)
		}

		buffer.WriteString(string(directory))
	}

	fmt.Fprintln(os.Stdout, buffer.String())
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: extend-path [-d|-a|-p] [-s] [-v] [-e VAR] DIR...")
	fmt.Fprintln(os.Stderr)
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "if none of delete, append or prepend is specified, extend-path")
	fmt.Fprintln(os.Stderr, "will only check the presence of the entries in PATH. if no")
	fmt.Fprintln(os.Stderr, "entries are specified, extend-path will simply print the current")
	fmt.Fprintln(os.Stderr, "PATH value to stdout.")
	fmt.Fprintln(os.Stderr, "each directory is converted to an absolute path using the")
	fmt.Fprintln(os.Stderr, "current directory if not already absolute.")
}
