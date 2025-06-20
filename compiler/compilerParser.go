/*
This file?
  - This is (or will be) an implementation of pratt parsing for operator precedence
  - This will mostly be paralellized.
*/
package compiler

import (
	"fmt"
	"os"
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

type TaskDep interface {
	Evaluate() error
	Exists() bool
}

type TaskDefInput struct {
	FileName string
	Exist    bool
}

func (t *TaskDefInput) Evaluate() error {
	if t.Exists() {
		return nil
	} else {
		return fmt.Errorf("file doesn't exist")
	}
}

// NOTE: doesn't allow inputs to be symlinks (might change in future)
func (t *TaskDefInput) Exists() bool {
	if _, err := os.Lstat(t.FileName); err == nil {
		return true
	}
	return false
}

type TaskDefOutput struct {
	FileName string
	Exists   bool
}

type TaskDef struct {
	Name    string
	Inputs  []*TaskDefInput
	Outputs []*TaskDefOutput
}
