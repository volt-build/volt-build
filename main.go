package main

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
	l "github.com/volt-build/volt-build/language"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Use all CPU cores

	var (
		silent     bool
		verbose    bool
		singleTask string
	)

	cmd := &cobra.Command{
		Use:     "volt-build [optional_path] [-s|--silent] [-v|--verbose] [-V|--version] [-t|--task <task>]",
		Short:   "A small build system focused on simplicity and speed.",
		Version: "0.1.1",
		Args:    cobra.MaximumNArgs(1),
		Long: `A small build system focused on simplicity and speed.
Supports incremental rebuilds. More features coming soon!`,
		Run: func(cmd *cobra.Command, args []string) {
			if silent && verbose {
				fmt.Fprintln(os.Stderr, "\x1b[1;31merror:\x1b[0m cannot mix silent and verbose flags")
				os.Exit(69)
			}

			// Default to ./build.volt if no path provided
			path := "./build.volt"
			if len(args) == 1 {
				path = args[0] + "/build.volt"
			}

			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\x1b[1;31merror:\x1b[0m %v\n try -h for help \n", err)
				os.Exit(1)
			}

			mode := getMode(silent, verbose)

			if singleTask != "" {
				// Run just one task from the build file
				if err := l.RunSingleTask(string(content), singleTask, mode); err != nil {
					fmt.Fprintf(os.Stderr, "\x1b[1;31merror:\x1b[0m %v\n", err)
					os.Exit(69)
				}
			} else {
				// Run the entire script
				if err := l.RunTaskScript(string(content), mode); err != nil {
					fmt.Fprintf(os.Stderr, "\x1b[1;31merror:\x1b[0m %v\n", err)
					os.Exit(1)
				}
			}
		},
	}

	cmd.CompletionOptions.DisableDefaultCmd = true

	// CLI flags
	cmd.Flags().StringVarP(&singleTask, "task", "t", "", "Run a single task from the build file")
	cmd.Flags().BoolVarP(&silent, "silent", "s", false, "Silent evaluation (no output)")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose evaluation (detailed output)")

	// Execute the command using fang
	if err := fang.Execute(context.Background(), cmd); err != nil {
		os.Exit(1)
	}
}

// Select evaluation mode based on flags
func getMode(silent, verbose bool) l.EvalMode {
	switch {
	case verbose:
		return l.EvalVerbose
	case silent:
		return l.EvalSilent
	default:
		return l.EvalRegular
	}
}
