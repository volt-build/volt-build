package main

import (
	"flag"
	"fmt"
	"os"
)

// FIXME: flags, please.
func main() {
	silent := flag.Bool("silent", false, "Don't print anything out and suppress all push statements")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	version := flag.Bool("version", false, "Print out the version and exit")
	task := flag.String("single-task", "", "To execute only a single task")
	flag.Parse()

	if *version && *silent {
		fmt.Fprintln(os.Stderr, "error: cannot mix 'version' and 'silent' flags")
		os.Exit(1)
	}

	if *version {
		fmt.Println("mini-build: version 0.1.0")
		return
	}

	runFromPath := func(path string, verbose, silent bool, task string) {
		input, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
		if task != "" {
			RunSingleTask(string(input), task)
		} else {
			RunTaskScript(string(input), verbose, silent)
		}
	}

	switch n := flag.NArg(); n {
	case 0:
		runFromPath("./Taskfile", *verbose, *silent, *task)
	case 1:
		path := flag.Arg(0) + "/Taskfile"
		runFromPath(path, *verbose, *silent, *task)
	default:
		fmt.Fprintln(os.Stderr, "cannot have more than 1 positional argument")
		os.Exit(1)
	}
}
