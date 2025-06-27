// only compiles if not windows (applies to everything. This program isn't desgined to work on windows, especially compilation.)
//go:build !windows

package codegen

import (
	"strings"

	l "github.com/volt-build/volt-build/language"
)

/*
- This type is to generate FASM (Flat Assembly) from a mini-build file.

	It is designed in the same way as the interpreter, it uses a strange
	technique of compilation that I like to call "recursive descent compilation"
	It works in a similar fashion to recursive descent parsing, where instead of
	the script being turned into an AST progressively, it's the AST being turned into
	target language. I will be implementing a faster method of compilation if this method
	seems to be too slow.
*/
type FasmGenerator struct {
	output *strings.Builder      // pointer to a strings.Builder to write to perodically output the generated assembly
	tasks  map[string]*l.TaskDef // very useful
	labels int                   // Amount of labels in the generated file (append over time)

	// current Label name (in fasm) (set to start when in entrypoint and set to data/rodata when declaring variables)
	currentLabel string
}

// Make a new FasmGenerator
func NewFasmGenerator(outputBuilder *strings.Builder, file []byte) *FasmGenerator {
	lexer := l.NewLexer(string(file))
	parser := l.NewParser(lexer)
	return &FasmGenerator{
		output: outputBuilder,
		tasks:  countTaskDefNodes(parser.ParseProgram()),
	}
}

// Count the task definitions in a program (AST Node/Leaf)
func countTaskDefNodes(program *l.Program) map[string]*l.TaskDef {
	returnMap := make(map[string]*l.TaskDef)
	for _, stmt := range program.Statements {
		if stmt.Type() == l.TaskDefNode {
			taskDef := stmt.(*l.TaskDef)
			returnMap[taskDef.Name] = taskDef
		}
	}
	return returnMap
}
