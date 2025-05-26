package main

import (
	"errors"
	"fmt"
	"log"
)

type EvalMode int

// there's a lot of enums in this project
// enum for mode evaluation
const (
	_ EvalMode = iota
	EvalSilent
	EvalVerbose
	EvalRegular
)

// Main way to interact with everything.
func RunTaskScript(input string, mode EvalMode) error {
	lex := NewLexer(input)
	parser := NewParser(lex)
	program := parser.ParseProgram()

	if len(parser.errors) > 0 {
		for _, err := range parser.errors {
			log.Printf("%v\n", err)
		}
		return errors.New("parsing failed")
	}
	interpreter := NewInterpreter()
	switch mode {
	case EvalRegular:
		_, err := interpreter.Evaluate(program)
		if err != nil {
			fmt.Printf("evaluation failed: %v\n", err)
			return err
		}
	case EvalVerbose:
		_, err := interpreter.EvaluateVerbosely(program)
		if err != nil {
			fmt.Printf("evaluation failed: %v\n", err)
			return err
		}
	case EvalSilent:
		_, err := interpreter.EvaluateWithoutPrinting(program)
		if err != nil {
			fmt.Printf("evaluation failed: %v\n", err)
			return err
		}
	default:
		fmt.Printf("Invalid evalMode.\n")
		return errors.New("invalid eval")
	}
	return nil
}

// here, input is the whole file and taskName is the name of the task to be run without variables
// Run a single task with this function for flags/position arguments;
func RunSingleTask(input string, taskName string) error {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.errors) > 0 {
		for _, err := range parser.errors {
			fmt.Printf("%v\n", err)
		}
		return errors.New("parsing failed")
	}

	interpreter := NewInterpreter()
	for _, stmt := range program.Statements {
		if stmt.Type() == TaskDefNode {
			task := stmt.(*TaskDef)
			interpreter.env.RegisterTask(task)
		}
	}

	task, exists := interpreter.env.GetTask(taskName)
	if !exists {
		return fmt.Errorf("task does not exist : %s", task)
	}

	execStmt := &ExecStatement{TaskName: taskName}
	_, err := interpreter.evaluateExec(execStmt)
	return err
}

func RunTaskScriptWithStatus(input string /* mode EvalMode // not yet supported */) error {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	progrem := parser.ParseProgram()

	if len(parser.errors) > 0 {
		for _, i := range parser.errors {
			fmt.Printf("%s\n", i)
		}
		return errors.New("parsing failed")
	}

	interpreter := NewInterpreter()
	_, err := interpreter.Evaluate(progrem)
	return err
}
