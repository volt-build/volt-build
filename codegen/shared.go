package codegen

import l "github.com/randomdude16671/mini-build/language"

// the target currently being compiled to.
type CompilerTarget int

// Append to over time
const (
	_          CompilerTarget = iota
	TargetELF                 // Compiles to flat assembly (FASM)
	TargetBEAM                // Compiles to erlang  (Not implemented.)
)

// Enum to hold types of symbols
type SymbolType int

const (
	_ SymbolType = iota
	// TODO: Finish this list
	TaskSymbol
	VariableSymbol
)

// Analyzer to put through before codegen
type SemanticAnalyzer struct {
	input         *string
	currentSymbol Symbol    // current symbol under analysis
	program       l.Program // program node to analyze
}

type Symbol interface {
	Type() SymbolType
	Finished() chan<- bool
}
