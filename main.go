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
		os.Exit(0)
	}

	var path string
	if flag.NArg() == 0 {
		path = "./Taskfile"
	} else if flag.NArg() == 1 {
		path = flag.Arg(0) + "/Taskfile"
	} else {
		fmt.Printf("Cannot have more than one positional argument\n")
		flag.Usage()
		os.Exit(1)
	}

	if *silent && *verbose {
		fmt.Printf("Cannot mix `silent` and `verbose` flags.\n")
		os.Exit(1)
	}

	input, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	if *silent {
		RunTaskScript(string(input), EvalSilent)
	} else if *verbose {
		RunTaskScript(string(input), EvalVerbose)
	} else {
		RunTaskScript(string(input), EvalRegular)
	}
}
