package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// NOTE: possibly will be increasing GOMAXPROCS

// `spawnCompile` is the name of this compile evaluation function because it uses workgroups and stuff
// with concurrency, so it runs on another OS thread
func (i *Interpreter) spawnCompile(cmpStmt *CompileStatement) (any, error) {
	type result struct {
		val any
		err error
	}
	resCh := make(chan result, 1)
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		fileExpr, err := i.Evaluate(cmpStmt.File)
		if err != nil {
			resCh <- result{nil, err}
			return
		}
		fileStr, ok := fileExpr.(string)
		if !ok {
			resCh <- result{nil, fmt.Errorf("compile file must evaluate to a string")}
			return
		}
		absolutePath, err := filepath.Abs(fileStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		cmdExpr, err := i.Evaluate(cmpStmt.Command)
		if err != nil {
			resCh <- result{nil, err}
			return
		}
		cmdStr, ok := cmdExpr.(string)
		if !ok {
			resCh <- result{nil, fmt.Errorf("compile command must evaluate to a string")}
			return
		}

		cmd := exec.Command("sh", "-c", cmdStr+" "+absolutePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			i.env.lastExitCode = 1
		} else {
			i.env.lastExitCode = 0
		}
		i.env.progressDone++
		resCh <- result{nil, err}
	}()
	wg.Wait()
	res := <-resCh
	return res.val, res.err
}

func (i *Interpreter) evaluateCompileWithoutPrinting(cmpStmt *CompileStatement) (any, error) {
	fileExpr, err := i.EvaluateWithoutPrinting(cmpStmt.File)
	if err != nil {
		return nil, err
	}

	fileStr, ok := fileExpr.(string)
	if !ok {
		return nil, fmt.Errorf("compile file must be a string, else its not used and throws this error")
	}

	cmdExpr, err := i.EvaluateWithoutPrinting(cmpStmt.Command)
	if err != nil {
		return nil, err
	}

	cmdStr, ok := cmdExpr.(string)
	if !ok {
		return nil, fmt.Errorf("compile command must be a string")
	}

	cmd := exec.Command("sh", "-c", cmdStr+" "+fileStr)
	// Don't redirect stdout/stderr to suppress output
	cmd.Stdout = nil
	cmd.Stderr = nil

	// initalize a channel to fill when command is done on seperate goroutine (cuz i like speed)
	errCh := make(chan error, 1)

	// run on seperate goroutine
	go func() {
		errCh <- cmd.Run()
	}()

	err = <-errCh
	if err != nil {
		i.env.lastExitCode = 1
	} else {
		i.env.lastExitCode = 0
	}

	i.env.progressDone++
	return nil, err
}

func (i *Interpreter) evaluateCompileVerbose(cmpStmt *CompileStatement) (any, error) {
	fileExpr, err := i.Evaluate(cmpStmt.File)
	if err != nil {
		return nil, err
	}

	fileStr, ok := fileExpr.(string)
	if !ok {
		return nil, fmt.Errorf("compile file must be a string, else its not used and throws this error")
	}

	cmdExpr, err := i.Evaluate(cmpStmt.Command)
	if err != nil {
		return nil, err
	}

	cmdStr, ok := cmdExpr.(string)
	if !ok {
		return nil, fmt.Errorf("compile command must be a string")
	}

	fmt.Printf("compiling: %s; with command: %s\n", cmpStmt.File, cmpStmt.Command)
	cmd := exec.Command("sh", "-c", cmdStr+" "+fileStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// initalize a channel to fill when command is done on seperate goroutine (cuz i like speed)
	errCh := make(chan error, 1)
	fmt.Printf("errCh initalized\n")

	// run on seperate goroutine
	go func() {
		errCh <- cmd.Run()
	}()
	fmt.Printf("errCh filled on different goroutine\n")
	err = <-errCh
	fmt.Printf("err: %v\n", err)
	if err != nil {
		i.env.lastExitCode = 1
		fmt.Printf("last exit code resulted in: %d\n", i.env.lastExitCode)
	} else {
		i.env.lastExitCode = 0
		fmt.Printf("last exit code resulted in: %d\n", i.env.lastExitCode)
	}

	return nil, err
}

func (i *Interpreter) evaluateConcat(concatOp *ConcatOperation) (any, error) {
	leftVal, err := i.Evaluate(concatOp.Left)
	if err != nil {
		return nil, err
	}

	rightVal, err := i.Evaluate(concatOp.Right)
	if err != nil {
		return nil, err
	}
	leftStr := fmt.Sprintf("%v", leftVal)
	rightStr := fmt.Sprintf("%v", rightVal)

	return leftStr + rightStr, nil
}

func (i *Interpreter) evaluateConcatWithoutPrinting(concatOp *ConcatOperation) (any, error) {
	leftVal, err := i.EvaluateWithoutPrinting(concatOp.Left)
	if err != nil {
		return nil, err
	}

	rightVal, err := i.EvaluateWithoutPrinting(concatOp.Right)
	if err != nil {
		return nil, err
	}
	leftStr := fmt.Sprintf("%v", leftVal)
	rightStr := fmt.Sprintf("%v", rightVal)

	return leftStr + rightStr, nil
}

func (i *Interpreter) evaluateProgram(p *Program) (any, error) {
	var result any
	var err error

	for _, stmt := range p.Statements {
		if stmt.Type() == TaskDefNode {
			task := stmt.(*TaskDef)
			i.env.RegisterTask(task)
		}
	}

	for _, stmt := range p.Statements {
		result, err = i.Evaluate(stmt)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (i *Interpreter) evaluateProgramWithoutPrinting(p *Program) (any, error) {
	var result any
	var err error

	for _, stmt := range p.Statements {
		if stmt.Type() == TaskDefNode {
			task := stmt.(*TaskDef)
			i.env.RegisterTask(task)
		}
	}

	for _, stmt := range p.Statements {
		result, err = i.EvaluateWithoutPrinting(stmt)
		if err != nil {
			return nil, err
		}

	}

	return result, nil
}

func (i *Interpreter) evaluateTaskDef(task *TaskDef) (any, error) {
	// Task definitions are handled in first pass.
	return nil, nil
}

func (i *Interpreter) evaluateTaskDefWithoutPrinting(task *TaskDef) (any, error) {
	// Task definitions are handled in first pass.
	return nil, nil
}

func (i *Interpreter) evaluateExec(execStmt *ExecStatement) (any, error) {
	task, exists := i.env.GetTask(execStmt.TaskName)
	if !exists {
		return nil, fmt.Errorf("task %s not found", execStmt.TaskName)
	}

	for _, depName := range task.Dependencies {
		dep, exists := i.env.GetTask(depName)
		if !exists {
			return nil, fmt.Errorf("task %s doesn't exist", depName)
		}

		_, err := i.Evaluate(dep.Body)
		if err != nil {
			return nil, err
		}
	}

	return i.Evaluate(task.Body)
}

func (i *Interpreter) evaluateExecWithoutPrinting(execStmt *ExecStatement) (any, error) {
	task, exists := i.env.GetTask(execStmt.TaskName)
	if !exists {
		return nil, fmt.Errorf("task %s not found", execStmt.TaskName)
	}

	for _, depName := range task.Dependencies {
		dep, exists := i.env.GetTask(depName)
		if !exists {
			return nil, fmt.Errorf("task %s doesn't exist", depName)
		}

		_, err := i.EvaluateWithoutPrinting(dep.Body)
		if err != nil {
			return nil, err
		}
	}

	return i.EvaluateWithoutPrinting(task.Body)
}

func (i *Interpreter) evaluateExecVerbose(execStmt *ExecStatement) (any, error) {
	task, exists := i.env.GetTask(execStmt.TaskName)
	fmt.Printf("found task %s, exists: %s", task, strconv.FormatBool(exists))
	if !exists {
		return nil, fmt.Errorf("task %s not found", execStmt.TaskName)
	}

	for _, depName := range task.Dependencies {
		dep, exists := i.env.GetTask(depName)
		if !exists {
			return nil, fmt.Errorf("task %s doesn't exist", depName)
		}

		_, err := i.Evaluate(dep.Body)
		if err != nil {
			return nil, err
		}
	}

	return i.Evaluate(task.Body)
}

func (i *Interpreter) evaluateShell(shellStmt *ShellStatement) (any, error) {
	cmdExpr, err := i.Evaluate(shellStmt.Command)
	if err != nil {
		return nil, err
	}

	cmdStr, ok := cmdExpr.(string)
	if !ok {
		return nil, fmt.Errorf("shell command must be a string")
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the command on different goroutine so its atleast a bit parallelized. (cuz its the start)
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Run()
	}()

	i.env.progressDone++
	return nil, <-errCh
}

func (i *Interpreter) evaluateShellWithoutPrinting(shellStmt *ShellStatement) (any, error) {
	cmdExpr, err := i.EvaluateWithoutPrinting(shellStmt.Command)
	if err != nil {
		return nil, err
	}

	cmdStr, ok := cmdExpr.(string)
	if !ok {
		return nil, fmt.Errorf("shell command must be a string")
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	// Don't redirect stdout/stderr to suppress output
	cmd.Stdout = nil
	cmd.Stderr = nil

	// run the command on different goroutine so its atleast a bit parallelized. (cuz its the start)
	errCh := make(chan error, 1)
	go func() {
		errCh <- cmd.Run()
	}()

	i.env.progressDone++
	return nil, <-errCh
}

func (i *Interpreter) evaluateShellVerbose(shellStmt *ShellStatement) (any, error) {
	cmdExpr, err := i.Evaluate(shellStmt.Command)
	if err != nil {
		return nil, err
	}

	cmdStr, ok := cmdExpr.(string)
	if !ok {
		return nil, fmt.Errorf("shell command must be a string")
	}

	cmd := exec.Command("sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the command on different goroutine so its atleast a bit parallelized. (cuz its the start)
	errCh := make(chan error, 1)
	fmt.Printf("Running command: sh -c %s\n", shellStmt.Command)
	go func() {
		errCh <- cmd.Run()
	}()
	fmt.Printf("error channel filled\n")
	fmt.Printf("returning\n")
	i.env.progressDone++
	return nil, <-errCh
}

func (i *Interpreter) evaluatePush(pushStmt *PushStatement) (any, error) {
	val, err := i.Evaluate(pushStmt.Value)
	if err != nil {
		return nil, err
	}
	fmt.Println(val) // Move print here
	return val, nil
}

func (i *Interpreter) evaluatePushWithoutPrinting(pushStmt *PushStatement) (any, error) {
	val, err := i.EvaluateWithoutPrinting(pushStmt.Value)
	if err != nil {
		return nil, err
	}
	i.env.SetVariable("PushStmtResult", val)
	return val, nil
}

func (i *Interpreter) evaluateIf(ifStmt *IfStatement) (any, error) {
	condition, err := i.Evaluate(ifStmt.Condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return i.Evaluate(ifStmt.ThenBlock)
	} else if ifStmt.ElseBlock != nil {
		return i.Evaluate(ifStmt.ElseBlock)
	}

	return nil, nil
}

func (i *Interpreter) evaluateIfWithoutPrinting(ifStmt *IfStatement) (any, error) {
	condition, err := i.EvaluateWithoutPrinting(ifStmt.Condition)
	if err != nil {
		return nil, err
	}

	if isTruthy(condition) {
		return i.EvaluateWithoutPrinting(ifStmt.ThenBlock)
	} else if ifStmt.ElseBlock != nil {
		return i.EvaluateWithoutPrinting(ifStmt.ElseBlock)
	}

	return nil, nil
}

func (i *Interpreter) evaluateForEach(forEachStmt *ForEachStatement) (any, error) {
	pattern := forEachStmt.Pattern

	if !strings.HasPrefix(pattern, "\"") && !strings.HasPrefix(pattern, "'") {
		if val, exists := i.env.GetVariable(pattern); exists {
			if str, ok := val.(string); ok {
				pattern = str
			} else {
				return nil, fmt.Errorf("foreach pattern must evaluate to a string, got %T", val)
			}
		}
	}

	if strings.HasPrefix(pattern, "\"") && strings.HasSuffix(pattern, "\"") {
		pattern = pattern[1 : len(pattern)-1]
	}
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var result any

	oldValue, exists := i.env.GetVariable(forEachStmt.VarName)
	for _, match := range matches {
		i.env.SetVariable(forEachStmt.VarName, match)

		result, err = i.Evaluate(forEachStmt.Body)
		if err != nil {
			if exists {
				i.env.SetVariable(forEachStmt.VarName, oldValue)
			} else {
				delete(i.env.variables, forEachStmt.VarName)
			}
		}
	}

	return result, nil
}

func (i *Interpreter) evaluateForEachWithoutPrinting(forEachStmt *ForEachStatement) (any, error) {
	pattern := forEachStmt.Pattern

	if !strings.HasPrefix(pattern, "\"") && !strings.HasPrefix(pattern, "'") {
		if val, exists := i.env.GetVariable(pattern); exists {
			if str, ok := val.(string); ok {
				pattern = str
			} else {
				return nil, fmt.Errorf("foreach pattern must evaluate to a string, got %T", val)
			}
		}
	}

	if strings.HasPrefix(pattern, "\"") && strings.HasSuffix(pattern, "\"") {
		pattern = pattern[1 : len(pattern)-1]
	}
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var result any

	oldValue, exists := i.env.GetVariable(forEachStmt.VarName)
	for _, match := range matches {
		i.env.SetVariable(forEachStmt.VarName, match)

		result, err = i.EvaluateWithoutPrinting(forEachStmt.Body)
		if err != nil {
			if exists {
				i.env.SetVariable(forEachStmt.VarName, oldValue)
			} else {
				delete(i.env.variables, forEachStmt.VarName)
			}
		}
	}

	return result, nil
}

func (i *Interpreter) evaluateForEachVerbose(forEachStmt *ForEachStatement) (any, error) {
	pattern := forEachStmt.Pattern
	fmt.Printf("foreach statement with %s as pattern\n", pattern)
	if !strings.HasPrefix(pattern, "\"") && !strings.HasPrefix(pattern, "'") {
		if val, exists := i.env.GetVariable(pattern); exists {
			if str, ok := val.(string); ok {
				fmt.Printf("%s\n", str)
				pattern = str
			} else {
				fmt.Printf("returning error..\n")
				return nil, fmt.Errorf("foreach pattern must evaluate to a string, got %T", val)
			}
		}
	}

	if strings.HasPrefix(pattern, "\"") && strings.HasSuffix(pattern, "\"") {
		pattern = pattern[1 : len(pattern)-1]
	}
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var result any

	oldValue, exists := i.env.GetVariable(forEachStmt.VarName)
	for _, match := range matches {
		i.env.SetVariable(forEachStmt.VarName, match)

		result, err = i.Evaluate(forEachStmt.Body)
		if err != nil {
			if exists {
				fmt.Printf("setting variable: %s, to value %T", forEachStmt.VarName, oldValue)
				i.env.SetVariable(forEachStmt.VarName, oldValue)
			} else {
				fmt.Printf("deleting variable: %s\n", forEachStmt.VarName)
				delete(i.env.variables, forEachStmt.VarName)
			}
		}
	}

	fmt.Printf("returning\n")
	return result, nil
}

func (i *Interpreter) evaluateBlock(blockStmt *BlockStatement) (any, error) {
	var result any
	var err error

	for _, stmt := range blockStmt.Statements {
		result, err = i.Evaluate(stmt)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (i *Interpreter) evaluateBlockWithoutPrinting(blockStmt *BlockStatement) (any, error) {
	var result any
	var err error

	for _, stmt := range blockStmt.Statements {
		result, err = i.EvaluateWithoutPrinting(stmt)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (i *Interpreter) evaluateBlockVerbose(blockStmt *BlockStatement) (any, error) {
	var result any
	var err error

	for _, stmt := range blockStmt.Statements {
		result, err = i.Evaluate(stmt)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%T\n", result)
	}
	fmt.Printf("returning\n")
	return result, nil
}

func (i *Interpreter) evaluateIdentifier(ident *Identifier) (any, error) {
	val, exists := i.env.GetVariable(ident.Value)
	if !exists {
		return nil, fmt.Errorf("variable named %s not found bro", ident.Value)
	}
	return val, nil
}

func (i *Interpreter) evaluateIdentifierWithoutPrinting(ident *Identifier) (any, error) {
	val, exists := i.env.GetVariable(ident.Value)
	if !exists {
		return nil, fmt.Errorf("variable named %s not found bro", ident.Value)
	}
	return val, nil
}

func (i *Interpreter) evaluateIdentifierVerbose(ident *Identifier) (any, error) {
	fmt.Printf("encoutered identifer, %s\n", ident.String())
	val, exists := i.env.GetVariable(ident.Value)
	if !exists {
		return nil, fmt.Errorf("variable named %s not found bro", ident.Value)
	}
	fmt.Printf("returning\n")
	return val, nil
}

func (i *Interpreter) evaluateShellExpr(shellExpr *ShellExpr) (any, error) {
	if shellExpr.Name == "?" {
		return i.env.lastExitCode, nil
	}

	// Fallback to actual OS environment
	if val, ok := os.LookupEnv(shellExpr.Name); ok {
		return val, nil
	}

	return nil, fmt.Errorf("unknown shell variable $%s", shellExpr.Name)
}

func (i *Interpreter) evaluateShellExprWithoutPrinting(shellExpr *ShellExpr) (any, error) {
	if shellExpr.Name == "?" {
		return i.env.lastExitCode, nil
	}

	// Fallback to actual OS environment
	if val, ok := os.LookupEnv(shellExpr.Name); ok {
		return val, nil
	}

	return nil, fmt.Errorf("unknown shell variable $%s", shellExpr.Name)
}

func (i *Interpreter) evaluateShellExprVerbose(shellExpr *ShellExpr) (any, error) {
	fmt.Printf("encountered shellExpr, %s\n", shellExpr.String())
	if shellExpr.Name == "?" {
		return i.env.lastExitCode, nil
	}

	// Fallback to actual OS environment
	if val, ok := os.LookupEnv(shellExpr.Name); ok {
		return val, nil
	}

	fmt.Printf("returning\n")
	return nil, fmt.Errorf("unknown shell variable $%s", shellExpr.Name)
}

func isTruthy(obj any) bool {
	switch v := obj.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	case nil:
		return false
	default:
		return true
	}
}

func (i *Interpreter) EvaluateWithoutPrinting(node Node) (any, error) {
	switch node.Type() {
	case ProgramNode:
		return i.evaluateProgramWithoutPrinting(node.(*Program))
	case TaskDefNode:
		i.env.RegisterTask(node.(*TaskDef))
		return i.evaluateTaskDefWithoutPrinting(node.(*TaskDef))
	case ExecNode:
		return i.evaluateExecWithoutPrinting(node.(*ExecStatement))
	case ShellNode:
		return i.evaluateShellWithoutPrinting(node.(*ShellStatement))
	case PushNode:
		return i.evaluatePushWithoutPrinting(node.(*PushStatement))
	case IfNode:
		return i.evaluateIfWithoutPrinting(node.(*IfStatement))
	case ForEachNode:
		return i.evaluateForEachWithoutPrinting(node.(*ForEachStatement))
	case BlockNode:
		return i.evaluateBlockWithoutPrinting(node.(*BlockStatement))
	case CompileNode:
		return i.evaluateCompileWithoutPrinting(node.(*CompileStatement))
	case StringNode:
		return node.(*StringLiteral).Value, nil
	case NumberNode:
		return node.(*NumberLiteral).Value, nil
	case IdentNode:
		return i.evaluateIdentifierWithoutPrinting(node.(*Identifier))
	case ShellExprNode:
		return i.evaluateShellExprWithoutPrinting(node.(*ShellExpr))
	case ConcatNode:
		return i.evaluateConcatWithoutPrinting(node.(*ConcatOperation))
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type())
	}
}

func (i *Interpreter) EvaluateVerbosely(node Node) (any, error) {
	switch node.Type() {
	case ProgramNode:
		return i.evaluateProgram(node.(*Program))
	case TaskDefNode:
		i.env.RegisterTask(node.(*TaskDef))
		return i.evaluateTaskDef(node.(*TaskDef))
	case ExecNode:
		return i.evaluateExecVerbose(node.(*ExecStatement))
	case ShellNode:
		return i.evaluateShellVerbose(node.(*ShellStatement))
	case PushNode:
		return i.evaluatePush(node.(*PushStatement))
	case IfNode:
		return i.evaluateIf(node.(*IfStatement))
	case ForEachNode:
		return i.evaluateForEachVerbose(node.(*ForEachStatement))
	case BlockNode:
		return i.evaluateBlockVerbose(node.(*BlockStatement))
	case CompileNode:
		return i.evaluateCompileVerbose(node.(*CompileStatement))
	case StringNode:
		return node.(*StringLiteral).Value, nil
	case NumberNode:
		return node.(*NumberLiteral).Value, nil
	case IdentNode:
		return i.evaluateIdentifierVerbose(node.(*Identifier))
	case ShellExprNode:
		return i.evaluateShellExprVerbose(node.(*ShellExpr))
	case ConcatNode:
		return i.evaluateConcat(node.(*ConcatOperation))
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type())
	}
}
