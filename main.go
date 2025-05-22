package main

import (
	"fmt"
	"os"
	// TODO: "github.com/fsnotify/fsnotify" -- File watching
)

func main() {
	bytes, err := os.ReadFile("./Taskfile")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	// Run the script itself
	RunTaskScript(string(bytes))
}
