package main

import (
	"errors"
	"fmt"
)

// Main way to interact with everything.
func RunTaskScript(input string, verbose bool, silent bool) error {
	// OPTIMIZE: consider replacing 2 args verbose and silent for only one enum (EvalMode)

	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.errors) > 0 {
		for _, err := range parser.errors {
			fmt.Println(err)
		}
		return errors.New("parsing failed")
	}

	interpreter := NewInterpreter()
	var err error
	if !verbose && !silent {
		_, err = interpreter.Evaluate(program)
	} else if verbose && !silent {
		_, err = interpreter.EvaluateVerbosely(program)
	} else if !verbose && silent {
		_, err = interpreter.EvaluateWithoutPrinting(program)
	} else {
		fmt.Printf("cannot mix flags verbose and silent\n")
		err = errors.New("conflicting flags")
	}

	return err
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
