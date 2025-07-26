package language

import (
	"fmt"
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
	INPUT:       "INPUT",
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
	Inputs       []string
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
