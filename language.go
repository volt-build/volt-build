package main

// another long file. Probably longer than the lexer.
// This file contains the AST
// Parser too.

import (
	"fmt"
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
	CompileNode    NodeType = "COMPILE"

	// Expression nodes
	IdentNode     NodeType = "IDENT"
	StringNode    NodeType = "STRING"
	NumberNode    NodeType = "NUMBER"
	BinaryOpNode  NodeType = "BINARY_OP"
	UnaryOpNode   NodeType = "UNARY_OP"
	ShellExprNode NodeType = "SHELL_EXPR"
	ConcatNode    NodeType = "CONCAT"
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
	CONCAT:      "++",
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
	SWAP:        "SWAP",
	WHILE:       "WHILE",
	FOREACH:     "FOREACH",
	COMPILE:     "COMPILE",
}

type Node interface {
	Type() NodeType
	String() string
}

type CompileStatement struct {
	File    Node
	Command Node
}

func (c *CompileStatement) Type() NodeType { return CompileNode }
func (c *CompileStatement) String() string {
	return fmt.Sprintf("compile %s %s", c.File.String(), c.Command.String())
}

type ConcatOperation struct {
	Left  Node
	Right Node
}

func (c *ConcatOperation) Type() NodeType { return ConcatNode }
func (c *ConcatOperation) String() string {
	return fmt.Sprintf("%s ++ %s", c.Left.String(), c.Right.String())
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
		if p.currentToken.Literal == "compile" {
			return p.parseCompileStatement()
		}
		return nil
	case COMPILE:
		return p.parseCompileStatement()
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

func (p *Parser) parseCompileStatement() *CompileStatement {
	stmt := &CompileStatement{}

	p.nextToken()                   // consume the "compile" keyword"
	stmt.File = p.parseExpression() // parse the expression in front of `compile`
	p.nextToken()                   // move onto the next expression

	stmt.Command = p.parseExpressionWithConcat()
	return stmt
}

func (p *Parser) parseExpressionWithConcat() Node {
	left := p.parseExpression()
	if p.peekTokenIs(CONCAT) {
		p.nextToken() // consume current token
		p.nextToken() // consume "compile"

		right := p.parseExpressionWithConcat()

		return &ConcatOperation{
			Left:  left,
			Right: right,
		}
	}

	return left
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

	p.nextToken() // consume `foreach`

	if !p.currentTokenIs(STRING) && !p.currentTokenIs(IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("Line %d, Column %d: expected string or identifier, got %s", p.currentToken.Line, p.currentToken.Column, p.currentToken.Literal))
		return nil
	}

	stmt.Pattern = p.currentToken.Literal
	stmt.VarName = "it"

	if p.peekTokenIs(IDENT) {
		p.nextToken()
		varName := p.currentToken.Literal
		if string(varName[0]) == "$" {
			varName = varName[1:]
		}
		stmt.VarName = varName
	}

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
	var left Node

	switch p.currentToken.Type {
	case STRING:
		left = &StringLiteral{Value: p.currentToken.Literal}
	case NUMBER:
		value, _ := strconv.ParseFloat(p.currentToken.Literal, 64)
		left = &NumberLiteral{Value: value}
	case IDENT:
		left = &Identifier{Value: p.currentToken.Literal}
	case SHELL:
		p.nextToken() // consume '$'
		left = &ShellExpr{Name: p.currentToken.Literal}
	default:
		return nil
	}

	// Check for concatenation operator
	if p.peekTokenIs(CONCAT) {
		p.nextToken() // consume the current token
		p.nextToken() // consume the CONCAT token

		right := p.parseExpression() // Handle right-associative concatenation

		return &ConcatOperation{
			Left:  left,
			Right: right,
		}
	}

	return left
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

func (i *Interpreter) GetTasks() map[string]*TaskDef {
	return i.env.tasks
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
	case CompileNode:
		return i.evaluateCompile(node.(*CompileStatement))
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
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type())
	}
}

func (i *Interpreter) EvaluateWithPrinting(node Node) (any, error) {
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
		return i.evaluatePushWithoutPrinting(node.(*PushStatement))
	case IfNode:
		return i.evaluateIf(node.(*IfStatement))
	case ForEachNode:
		return i.evaluateForEach(node.(*ForEachStatement))
	case BlockNode:
		return i.evaluateBlock(node.(*BlockStatement))
	case CompileNode:
		return i.evaluateCompile(node.(*CompileStatement))
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
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.Type())
	}
}
