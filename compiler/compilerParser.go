/*
This file?
  - This is (or will be) an implementation of pratt parsing for operator precedence
  - This will mostly be paralellized.
*/
package compiler

import (
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
	Ident     NodeType = "IDENT"
	String    NodeType = "STRING"
	Number    NodeType = "NUMBER"
	BinaryOp  NodeType = "BINARY_OP"
	UnaryOp   NodeType = "UNARY_OP"
	ShellExpr NodeType = "SHELL_EXPR"
	Concat    NodeType = "CONCAT"
)

type Node interface {
	Type() NodeType
	String() string
}

var lexerMap map[TokenType]string = map[TokenType]string{
	EOF:         "EOF",
	ILLEGAL:     "ILLEGAL",
	IDENT:       "IDENTIFIER",
	STRING:      "string",
	NEWLINE:     "newline",
	NUMBER:      "number",
	COMMENT:     "comment",
	DEFINE:      "define",
	ASSIGN:      "assign",
	EQUAL:       "equal",
	MINUS:       "minus",
	ASTERISK:    "asterisk",
	PLUS:        "plus",
	MODULO:      "modulo",
	SLASH:       "slash",
	GREATERTHAN: "greaterthan",
	LESSTHAN:    "lessthan",
	LORETO:      "loreto",
	GORETO:      "goreto",
	NOTEQUAL:    "notequal",
	NOT:         "not",
	AND:         "and",
	OR:          "or",
	SHELL:       "shell",
	CONCAT:      "concat",
	COMMA:       "comma",
	SEMICOLON:   "semicolon",
	LPAREN:      "lparen",
	RPAREN:      "rparen",
	LBRACE:      "lbrace",
	RBRACE:      "rbrace",
	LBRACKET:    "lbracket",
	RBRACKET:    "rbracket",
	PIPE:        "pipe",
	TASK:        "task",
	RUN:         "run",
	IF:          "if",
	ELSE:        "else",
	IMPORT:      "import",
	DEPENDENCY:  "dependency",
	SWAP:        "swap",
	WHILE:       "while",
	FOREACH:     "foreach",
	COMPILE:     "compile",
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

type TaskDefInput struct {
	Path     string   // Absolute path
	Optional bool     // self explanatory
	Hash     string   // Content hash for caching
	FromTask *TaskDef // A return from a task
}

type TaskDefOutput struct {
	Path     string // Absolute path
	Optional bool   // self explanatory
	Hash     string // Content hash for caching
}

type TaskDef struct {
	Name     string           // Name of the task
	TaskDeps []*TaskDef       // List of tasks to complete before this one
	Inputs   []*TaskDefInput  // List of inputs of the task
	Outputs  []*TaskDefOutput // List of outputs of the task
}
