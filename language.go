package main

// another long file. Probably longer than the lexer.
// This file contains the AST
// Parser too.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type NodeType string

const (
	// program node.
	ProgramNode NodeType = "PROGRAM"

	// statement nodes.
	TaskDefNode    NodeType = "TASK_DEF"
	ExecNode       NodeType = "EXEC"
	ShellNode      NodeType = "SHELL"
	PushNode       NodeType = "PUSH"
	IfNode         NodeType = "IF"
	ForEachNode    NodeType = "FOREACH"
	WhileNode      NodeType = "WHILE"
	ForNode        NodeType = "FOR"
	AssignmentNode NodeType = "ASSIGNMENT"
	BlockNode      NodeType = "BLOCK"

	// Expression nodes
	IdentNode     NodeType = "IDENT"
	StringNode    NodeType = "STRING"
	NumberNode    NodeType = "NUMBER"
	BinaryOpNode  NodeType = "BINARY_OP"
	UnaryOpNode   NodeType = "UNARY_OP"
	ShellExprNode NodeType = "SHELL_EXPR"
)

var lexerMap map[TokenType]string = map[TokenType]string{
	STRING:      "STR",
	IDENT:       "IDENT",
	NEWLINE:     "NEWLINE",
	NUMBER:      "NUMBER",
	COMMENT:     "COMMENT",
	DEFINE:      "define",
	ASSIGN:      "assign",
	EQUAL:       "equal",
	MINUS:       "minus",
	ASTERISK:    "asterisk",
	PLUS:        "plus",
	MODULO:      "modulo",
	SLASH:       "slash",
	GREATERTHAN: ">",
	LESSTHAN:    "<",
	LORETO:      "<=",
	GORETO:      ">=",
	NOTEQUAL:    "!=",
	NOT:         "!",
	AND:         "&&",
	OR:          "||",
	SHELL:       "$",
	COMMA:       "comma",
	SEMICOLON:   ";",
	LPAREN:      "(",
	RPAREN:      ")",
	LBRACE:      "{",
	RBRACE:      "}",
	LBRACKET:    "[",
	RBRACKET:    "]",
	PIPE:        "|",
	TASK:        "TASK",
	RUN:         "RUN",
	IF:          "IF",
	ELSE:        "ELSE",
	IMPORT:      "IMPORT",
	DEPENDENCY:  "DEPENDENCY",
	SWAP:        "SWA",
	LANGUAGE:    "LANGUAG",
	FOR:         "FOR",
	WHILE:       "WHILE",
	FOREACH:     "FOREACH",
}

type Node interface {
	Type() NodeType
	String() string
}

type Program struct {
	Statements []Node
}

func (p *Program) Type() NodeType { return ProgramNode }
func (p *Program) String() string {
	var out strings.Builder
	for _, stmt := range p.Statements {
		out.WriteString(stmt.String() + "\n")
	}
	return out.String()
}

type TaskDef struct {
	Name         string
	Dependencies []string
	Body         Node
}

func (t *TaskDef) Type() NodeType { return TaskDefNode }
func (t *TaskDef) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("task %s", t.Name))

	if len(t.Dependencies) > 0 {
		out.WriteString(" requires ")
		out.WriteString(strings.Join(t.Dependencies, ", "))
	}
	out.WriteString(" " + t.Body.String())
	return out.String()
}

type ExecStatement struct {
	TaskName string
}

func (e *ExecStatement) Type() NodeType { return ExecNode }
func (e *ExecStatement) String() string {
	return fmt.Sprintf("exec %s", e.TaskName)
}

type ShellStatement struct {
	Command Node
}

func (s *ShellStatement) Type() NodeType { return ShellNode }
func (s *ShellStatement) String() string {
	return fmt.Sprintf("shell %s", s.Command.String())
}

type PushStatement struct {
	Value Node // expr.
}

func (p *PushStatement) Type() NodeType { return PushNode }
func (p *PushStatement) String() string {
	return fmt.Sprintf("push %s", p.Value.String())
}

type IfStatement struct {
	Condition Node
	ThenBlock Node
	ElseBlock Node
}

func (i *IfStatement) Type() NodeType { return IfNode }
func (i *IfStatement) String() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("if %s %s", i.Condition.String(), i.ThenBlock.String()))
	if i.ElseBlock != nil {
		out.WriteString(fmt.Sprintf(" else %s", i.ElseBlock.String()))
	}
	return out.String()
}

type ForEachStatement struct {
	Pattern string
	VarName string
	Body    Node
}

func (f *ForEachStatement) Type() NodeType { return ForEachNode }
func (f *ForEachStatement) String() string {
	return fmt.Sprintf("foreach %s %s", f.Pattern, f.Body.String())
}

type BlockStatement struct {
	Statements []Node
}

func (b *BlockStatement) Type() NodeType { return BlockNode }
func (b *BlockStatement) String() string {
	var out strings.Builder
	out.WriteString("{\n")
	for _, stmt := range b.Statements {
		out.WriteString("  " + stmt.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

type AssignmentStatement struct {
	Name  string
	Value Node // Expression.
}

func (a *AssignmentStatement) Type() NodeType { return AssignmentNode }
func (a *AssignmentStatement) String() string {
	return fmt.Sprintf("%s = %s", a.Name, a.Value)
}

type Identifier struct {
	Value string
}

func (i *Identifier) Type() NodeType { return IdentNode }
func (i *Identifier) String() string { return i.Value }

type StringLiteral struct {
	Value string
}

func (s *StringLiteral) Type() NodeType { return StringNode }
func (s *StringLiteral) String() string {
	return fmt.Sprintf("\"%s\"", s.Value)
}

type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) Type() NodeType { return NumberNode }
func (n *NumberLiteral) String() string {
	return fmt.Sprintf("%g", n.Value)
}

type BinaryOperation struct {
	Left     Node
	Operator string
	Right    Node
}

func (b *BinaryOperation) Type() NodeType { return BinaryOpNode }
func (b *BinaryOperation) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left, b.Operator, b.Right)
}

type UnaryOperation struct {
	Operator string
	Operand  Node // Same like Right but makes more sense to call Operand since is the only thing.
}

func (u *UnaryOperation) Type() NodeType { return UnaryOpNode }
func (u *UnaryOperation) String() string {
	return fmt.Sprintf("(%s%s)", u.Operator, u.Operand)
}

type ShellExpr struct {
	Name string
}

func (s *ShellExpr) Type() NodeType { return ShellExprNode }
func (s *ShellExpr) String() string {
	return "$" + s.Name
}

type Parser struct {
	l            *Lexer // I WROTE IT!!  lol
	currentToken Token
	peekToken    Token
	errors       []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// To init peek and current token.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// no newlines and comments (skipping here)
	for p.peekToken.Type == NEWLINE || p.peekToken.Type == COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("Line %d, Column %d: expected next token to be %v, got %v instead", p.peekToken.Line, p.peekToken.Column, lexerMap[t], lexerMap[p.peekToken.Type])
	p.errors = append(p.errors, msg)
}

func (pr *Parser) ParseProgram() *Program {
	program := &Program{
		Statements: []Node{},
	}

	for !pr.currentTokenIs(EOF) {
		stmt := pr.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		pr.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Node {
	switch p.currentToken.Type {
	case TASK:
		return p.parseTaskDefinition()
	case IDENT:
		if p.peekTokenIs(ASSIGN) {
			return p.parseAssignStatement()
		}

		if p.currentToken.Literal == "exec" {
			return p.parseExecStatement()
		}
		if p.currentToken.Literal == "shell" {
			return p.parseShellStatement()
		}
		if p.currentToken.Literal == "push" {
			return p.parsePushStatement()
		}
		return nil
	case IF:
		return p.parseIfStatement()
	case FOREACH:
		return p.parseForEachStatement()
	case SHELL:
		return p.parseShellStatement()

	default:
		return nil
	}
}

func (p *Parser) parseTaskDefinition() *TaskDef {
	task := &TaskDef{}
	task.Dependencies = []string{} // init a empty slice for now.

	if !p.expectPeek(IDENT) {
		return nil
	}

	task.Name = p.currentToken.Literal

	if p.peekTokenIs(IDENT) && p.peekToken.Literal == "requires" {
		p.nextToken() // consume "requires"

		if !p.expectPeek(IDENT) {
			return nil
		}

		task.Dependencies = append(task.Dependencies, p.currentToken.Literal)

		for p.peekTokenIs(COMMA) {
			p.nextToken()
			if !p.expectPeek(IDENT) {
				return nil
			}
			task.Dependencies = append(task.Dependencies, p.currentToken.Literal)
		}
	}
	if !p.expectPeek(LBRACE) {
		return nil
	}

	task.Body = p.parseBlockStatement()
	return task
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{
		Statements: []Node{},
	}

	p.nextToken() // comsume `{`

	for !p.currentTokenIs(RBRACE) && !p.currentTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExecStatement() *ExecStatement {
	stmt := &ExecStatement{}
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.TaskName = p.currentToken.Literal
	return stmt
}

func (p *Parser) parseShellStatement() *ShellStatement {
	stmt := &ShellStatement{}
	p.nextToken()
	stmt.Command = p.parseExpression()

	return stmt
}

func (p *Parser) parsePushStatement() *PushStatement {
	stmt := &PushStatement{}
	p.nextToken()
	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{}

	p.nextToken() // consume if
	stmt.Condition = p.parseExpression()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.ThenBlock = p.parseBlockStatement()

	if p.peekTokenIs(ELSE) {
		p.nextToken() // consume else

		if !p.expectPeek(LBRACE) {
			return nil
		}
		stmt.ElseBlock = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseForEachStatement() *ForEachStatement {
	stmt := &ForEachStatement{}

	p.nextToken()

	if !p.currentTokenIs(STRING) && !p.currentTokenIs(IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("Line %d, Column %d: expected string or identifier, got %TokenType", p.currentToken.Line, p.currentToken.Column, p.currentToken.Type))
		return nil
	}

	stmt.Pattern = p.currentToken.Literal
	stmt.VarName = "$_"

	if !p.expectPeek(LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseAssignStatement() *AssignmentStatement {
	stmt := &AssignmentStatement{
		Name: p.currentToken.Literal,
	}

	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseExpression() Node {
	switch p.currentToken.Type {
	case STRING:
		return &StringLiteral{Value: p.currentToken.Literal}
	case NUMBER:
		value, _ := strconv.ParseFloat(p.currentToken.Literal, 64)
		return &NumberLiteral{Value: value}
	case IDENT:
		return &Identifier{Value: p.currentToken.Literal}
	case SHELL:
		p.nextToken() // consume '$'
		return &ShellExpr{Name: p.currentToken.Literal}
	default:
		return nil
	}
}

type Environment struct {
	variables    map[string]any
	tasks        map[string]*TaskDef
	lastExitCode int // store $? value in this
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

func (i *Interpreter) Evaluate(node Node) (any, error) {
	switch node.Type() {
	case ProgramNode:
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
	case StringNode:
		return node.(*StringLiteral).Value, nil
	case NumberNode:
		return node.(*NumberLiteral).Value, nil
	case IdentNode:
		return i.evaluateIdentifier(node.(*Identifier))
	case ShellExprNode:
		return i.evaluateShellExpr(node.(*ShellExpr))
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type())
	}
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

func (i *Interpreter) evaluateTaskDef(task *TaskDef) (any, error) {
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

func (i *Interpreter) evaluateIf(ifStmt *IfStatement) (any, error) {
	condition, err := i.Evaluate(ifStmt.ThenBlock)
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

func (i *Interpreter) evaluateForEach(forEachStmt *ForEachStatement) (any, error) {
	pattern := forEachStmt.Pattern

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var result any

	for _, match := range matches {
		i.env.SetVariable("_", match)

		result, err = i.Evaluate(forEachStmt.Body)
		if err != nil {
			return nil, err
		}
	}

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

func (i *Interpreter) evaluateIdentifier(ident *Identifier) (any, error) {
	val, exists := i.env.GetVariable(ident.Value)
	if !exists {
		return nil, fmt.Errorf("variable named %s not found bro", ident.Value)
	}
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
