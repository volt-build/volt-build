package language

// this file contains utils and the main run.

import (
	"fmt"
)

// TODO: add progress feedback
type Environment struct {
	variables     map[string]any      // variables inside the script
	tasks         map[string]*TaskDef // tasks to be executed
	progressDone  int                 // increment after all the compile/shell statements
	progressTotal int                 // total needed to be done
	lastExitCode  int                 // store $? value in this
}

func NewEnvironment() *Environment {
	return &Environment{
		variables:    make(map[string]any),
		tasks:        make(map[string]*TaskDef),
		lastExitCode: 0,
	}
}

func (env *Environment) SetVariable(name string, value any) {
	env.variables[name] = value
}

func (env *Environment) GetVariable(name string) (any, bool) {
	value, exists := env.variables[name]
	return value, exists
}

func (env *Environment) RegisterTask(task *TaskDef) {
	env.tasks[task.Name] = task
}

func (env *Environment) GetTask(name string) (*TaskDef, bool) {
	value, exists := env.tasks[name]
	return value, exists
}

type Interpreter struct {
	env *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnvironment(),
	}
}

func (i *Interpreter) GetTasks() map[string]*TaskDef {
	return i.env.tasks
}

func (i *Interpreter) Evaluate(node Node) (any, error) {
	switch node.Type() {
	case ProgramNode:
		i.preprocessEvaluateProgram(node.(*Program)) // increment counters for status.
		return i.evaluateProgram(node.(*Program))
	case TaskDefNode:
		i.env.RegisterTask(node.(*TaskDef))
		return i.evaluateTaskDef(node.(*TaskDef))
	case ExecNode:
		return i.evaluateExec(node.(*ExecStatement))
	case ShellNode:
		return i.evaluateShell(node.(*ShellStatement))
	case PushNode:
		return i.evaluatePush(node.(*PushStatement))
	case IfNode:
		return i.evaluateIf(node.(*IfStatement))
	case ForEachNode:
		return i.evaluateForEach(node.(*ForEachStatement))
	case BlockNode:
		return i.evaluateBlock(node.(*BlockStatement))
	case CompileNode:
		return i.spawnCompile(node.(*CompileStatement))
	case StringNode:
		return node.(*StringLiteral).Value, nil
	case NumberNode:
		return node.(*NumberLiteral).Value, nil
	case IdentNode:
		return i.evaluateIdentifier(node.(*Identifier))
	case ShellExprNode:
		return i.evaluateShellExpr(node.(*ShellExpr))
	case ConcatNode:
		return i.evaluateConcat(node.(*ConcatOperation))
	case AssignmentNode:
		return i.evaluateAssign(node.(*AssignmentStatement))
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type())
	}
}

func (i *Interpreter) evaluateAssign(assignStmt *AssignmentStatement) (any, error) {
	result, err := i.Evaluate(assignStmt.Value)
	if err != nil {
		return nil, err
	}
	i.env.SetVariable(assignStmt.Name, result)
	return result, nil
}

func (i *Interpreter) preprocessEvaluateProgram(p *Program) {
	i.env.progressTotal = 0 // Reset counter
	i.countExecutableStatements(p)
}

func (i *Interpreter) countExecutableStatements(node Node) {
	switch node.Type() {
	case ProgramNode:
		program := node.(*Program)
		for _, stmt := range program.Statements {
			i.countExecutableStatements(stmt)
		}
	case TaskDefNode:
		task := node.(*TaskDef)
		i.countExecutableStatements(task.Body)
	case BlockNode:
		block := node.(*BlockStatement)
		for _, stmt := range block.Statements {
			i.countExecutableStatements(stmt)
		}
	case CompileNode:
		i.env.progressTotal++
	case ShellNode:
		i.env.progressTotal++
	case IfNode:
		ifStmt := node.(*IfStatement)
		i.countExecutableStatements(ifStmt.ThenBlock)
		if ifStmt.ElseBlock != nil {
			i.countExecutableStatements(ifStmt.ElseBlock)
		}
	case ForEachNode:
		forEach := node.(*ForEachStatement)
		i.countExecutableStatements(forEach.Body)
	case ExecNode:
		exec := node.(*ExecStatement)
		if task, exists := i.env.GetTask(exec.TaskName); exists {
			i.countExecutableStatements(task.Body)
		}
	default:
	}
}
