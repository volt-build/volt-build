package main

import "fmt"

type VarType int

const (
	_ VarType = iota
	CompInt
	CompFloat
	CompBool
	CompString
)

type CompVar interface {
	Type() VarType // Return the type of the variable
	Value() any    // Return the value of the variable
}

type IntVar struct {
	Val    int    // return the value as an int (obviously)
	String string // Wrap the value in a string
}

func (i *IntVar) Type() VarType { return CompInt }
func (i *IntVar) Value() string {
	return fmt.Sprintf("%d", i.Val)
}

type FloatVar struct {
	Val    float64
	String string
}

func (f *FloatVar) Type() VarType { return CompFloat }
func (f *FloatVar) Value() string {
	return fmt.Sprintf("%f", f.Val)
}

type CompilerEnv struct {
	LastExitCode int
	Variables    map[CompVar]any
}

type Compiler struct{}
