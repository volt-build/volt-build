package language

import (
	"errors"
	"fmt"
	"log"
	"os"
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

func Exists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true // different error, doesn't need to be detected
		}
	}
	return true
}

func RunTaskScript(input string, mode EvalMode) error {
	err := os.MkdirAll("./.volt-build/", 0o755)
	if err != nil {
		return err
	}

	gitignoreFile, err := os.Create("./.volt-build/.gitignore")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	_, err = gitignoreFile.WriteString("*")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}

	gitignoreFile.Close()

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
	timestamps, err := interpreter.loadTimestamps(TIMESTAMP_PATH)
	if err != nil {
		fmt.Printf("Failed to load timestamps for incremental rebuilds: %v\n", err)
	}
	interpreter.timestamps = timestamps

	switch mode {
	case EvalRegular:
		_, err = interpreter.Evaluate(program)
		if err != nil {
			fmt.Printf("evaluation failed: %v\n", err)
			return err
		}
	case EvalVerbose:
		_, err = interpreter.EvaluateVerbosely(program)
		if err != nil {
			fmt.Printf("evaluation failed: %v\n", err)
			return err
		}
	case EvalSilent:
		_, err = interpreter.EvaluateWithoutPrinting(program)
		if err != nil {
			fmt.Printf("evaluation failed: %v\n", err)
			return err
		}
	default:
		fmt.Printf("Invalid evalMode.\n")
		return errors.New("invalid eval")
	}

	// Save timestamps after successful execution
	err = interpreter.saveTimestamps(TIMESTAMP_PATH, interpreter.timestamps)
	if err != nil {
		return err
	}
	return nil
}

// Updated RunSingleTask function
func RunSingleTask(input string, taskName string, mode EvalMode) error {
	err := os.MkdirAll("./.volt-build/", 0o755)
	if err != nil {
		return fmt.Errorf("failed to create .volt-build directory: %w", err)
	}

	gitignoreFile, err := os.Create("./.volt-build/.gitignore")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}
	_, err = gitignoreFile.WriteString("*")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return err
	}

	gitignoreFile.Close()

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
	timestamps, err := interpreter.loadTimestamps(TIMESTAMP_PATH)
	if err != nil {
		fmt.Printf("Failed to load timestamps for incremental rebuilds: %v\n", err)
	}
	interpreter.timestamps = timestamps

	// Register tasks
	for _, stmt := range program.Statements {
		if stmt.Type() == TaskDefNode {
			task := stmt.(*TaskDef)
			interpreter.env.RegisterTask(task)
		}
	}

	task, exists := interpreter.env.GetTask(taskName)
	if !exists {
		return fmt.Errorf("task does not exist: %s", taskName)
	}

	execStmt := &ExecStatement{TaskName: task.Name}

	switch mode {
	case EvalSilent:
		_, err = interpreter.evaluateExecWithoutPrinting(execStmt)
	case EvalVerbose:
		_, err = interpreter.evaluateExecVerbose(execStmt)
	case EvalRegular:
		_, err = interpreter.evaluateExec(execStmt)
	default:
		return fmt.Errorf("invalid EvalMode")
	}

	if err != nil {
		return err
	}

	// Save timestamps after execution
	if saveErr := interpreter.saveTimestamps(TIMESTAMP_PATH, interpreter.timestamps); saveErr != nil {
		return fmt.Errorf("failed to save timestamps: %w", saveErr)
	}

	return nil
}
