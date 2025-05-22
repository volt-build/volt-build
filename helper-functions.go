package main

import (
	"errors"
	"fmt"
)

// Main way to interact with everything.
func RunTaskScript(input string) error {
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
	_, err := interpreter.Evaluate(program)
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

// Run a single task with its variables in a map
// TODO: implement a way to just parse the variables into a map automatically;
func RunSingleTaskWithVariables(input string, taskName string, vars map[string]any) error {
	lexer := NewLexer(input)
	parser := NewParser(lexer)
	program := parser.ParseProgram()

	if len(parser.errors) > 0 {
		for _, err := range parser.errors {
			fmt.Printf("%s\n", err)
		}
		return errors.New("paring failed")
	}

	interpreter := NewInterpreter()

	for name, value := range vars {
		interpreter.env.SetVariable(name, value)
	}

	for _, stmt := range program.Statements {
		if stmt.Type() == TaskDefNode {
			task := stmt.(*TaskDef)
			interpreter.env.RegisterTask(task)
		}
	}

	_, exists := interpreter.env.GetTask(taskName)
	if !exists {
		return fmt.Errorf("task %s not found", taskName)
	}
	execStmt := &ExecStatement{TaskName: taskName}
	_, err := interpreter.evaluateExec(execStmt)
	return err
}
