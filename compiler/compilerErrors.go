package compiler

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func returnPointersToFile() {
}

func getLineFromFile(path string, lineNo int) (string, error) {
	var out strings.Builder
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	current := 1
	for scanner.Scan() {
		if current >= lineNo-1 && current <= lineNo+1 {
			out.WriteString(scanner.Text())
			out.WriteByte('\n')
		}
		if current > lineNo+1 {
			break
		}
		current++
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return out.String(), nil
}

type InputMissingError struct {
	File         string
	NotFoundFile string
	Op           string
	Err          error

	// for better error messages with nice looking diagnostics
	Line int
	Col  int
}

func (im *InputMissingError) Error() string {
	return fmt.Sprintf("%s: %v", im.Op, im.Err)
}

// The desgin is very human
func (im *InputMissingError) Human() string {
	var out strings.Builder
	out.WriteString("\n")
	out.WriteString("error: Input is missing:\n")
	lines, err := getLineFromFile(im.File, im.Line)

	return out.String()
}
