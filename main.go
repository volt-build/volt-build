package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

// pathExists checks if the given path exists.
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func main() {
	// Define flags
	silent := flag.Bool("silent", false, "Don't put anything to stdout, unless there is a push statement in the build script.")
	version := flag.Bool("version", false, "Print version and exit.")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] [script-dir]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	// Determine script directory (default: current directory)
	scriptDir := "."
	if flag.NArg() >= 1 {
		scriptDir = flag.Arg(0)
	}

	taskfilePath := scriptDir + "/Taskfile"

	if *version {
		fmt.Printf("mini-build version: 0.1.0\n")
		fmt.Printf("Exiting with code: 1")
		os.Exit(1)
	}

	if !*silent {
		if pathExists(scriptDir) {
			fmt.Printf("Running Task Script: %s\n", taskfilePath)
		} else {
			fmt.Printf("Path not found: %s\n", scriptDir)
		}
	} else {
		if pathExists(taskfilePath) {
			input, err := os.ReadFile(taskfilePath)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", taskfilePath, err)
				os.Exit(1)
			}
			RunTaskScript(string(input))
		} else {
			fmt.Printf("Path does not exist: %s\n", taskfilePath)
			os.Exit(1)
		}
	}
}
