/*
This file?
  - This is (or will be) an implementation of pratt parsing for operator precedence
  - This will mostly be paralellized.
*/
package compiler

type NodeType string

const (
	// program node.
	Program NodeType = "PROGRAM"

	// statement nodes.
	TaskDef    NodeType = "TASK_DEF"
	Exec       NodeType = "EXEC"
	Shell      NodeType = "SHELL"
	Push       NodeType = "PUSH"
	If         NodeType = "IF"
	ForEach    NodeType = "FOREACH"
	While      NodeType = "WHILE"
	For        NodeType = "FOR"
	Assignment NodeType = "ASSIGNMENT"
	Block      NodeType = "BLOCK"
	Compile    NodeType = "COMPILE"

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

var lexerMap map[TokenType]string = map[TokenType]string{}
