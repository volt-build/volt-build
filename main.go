package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	silent := flag.Bool("silent", false, "Enable silent mode")
	verbose := flag.Bool("verbose", false, "Enable verbose mode")
	version := flag.Bool("version", false, "Print the version and exit")
	singleTask := flag.String("t", "", "To use only a single task")

	flag.Parse()

	if *version {
		fmt.Println("mini-build: version 0.1.0")
		return
	}

	// Validate mutually exclusive flags
	if *silent && *verbose {
		fmt.Println("Error: can't mix -silent and -verbose")
		os.Exit(1)
	}

	// Determine evaluation mode
	var mode EvalMode
	switch {
	case *silent:
		mode = EvalSilent
	case *verbose:
		mode = EvalVerbose
	default:
		mode = EvalRegular
	}

	// Handle single task if specified
	if *singleTask != "" {
		input, err := os.ReadFile("./Taskfile")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		err = RunSingleTask(string(input), *singleTask, mode)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Handle full taskfile
	content, err := os.ReadFile("./Taskfile")
	if err != nil {
		fmt.Printf("Error reading Taskfile: %v\n", err)
		os.Exit(1)
	}

	if err := RunTaskScript(string(content), mode); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
