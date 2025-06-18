/*
This file?
  - This is (or will be) an implementation of pratt parsing for operator precedence
  - This will mostly be paralellized.
*/
package compiler

type NodeType string

const (
// TODO: define node types
)

type Node interface {
	Type() NodeType
	String() string
}
