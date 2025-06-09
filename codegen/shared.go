package codegen

type CompilerTarget int

const (
	_          CompilerTarget = iota
	TargetELF                 // Compiles to flat assembly (FASM)
	TargetBEAM                // Compiles to erlang  (Not implemented.)
)
