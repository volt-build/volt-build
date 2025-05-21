package main

import (
	"fmt"
	"os"
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
