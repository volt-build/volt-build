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

	flag.Parse()

	if *version {
		fmt.Printf("mini-build: version 0.1.0\n")
	}
	if *silent && *verbose {
		fmt.Printf("Can't mix those flags\n")
		return
	}
	if *silent {
		input, err := os.ReadFile("./Taskfile")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		RunTaskScript(string(input), EvalSilent)
	}
	if *verbose {
		input, err := os.ReadFile("./Taskfile")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		RunTaskScript(string(input), EvalVerbose)
	}
	if !*verbose && !*silent {
		input, err := os.ReadFile("./Taskfile")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		RunTaskScript(string(input), EvalRegular)
	}
}
