// Package object holds the values a running Pharos program works with.
package object

import (
	"fmt"

	"github.com/abhinavtripathy/pharos/interpreter/ast"
)

// Type names deliberately match the declaration keywords, so that an error
// about a value reads in the same words the program was written in.
type Type string

const (
	NumType       Type = "num"
	StringType    Type = "string"
	BoolType      Type = "bool"
	NothingType   Type = "nothing"
	CodeblockType Type = "codeblock"

	returnType Type = "give back"
	errorType  Type = "error"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Number struct{ Value float64 }

func (n *Number) Type() Type      { return NumType }
func (n *Number) Inspect() string { return ast.FormatNumber(n.Value) }

type Text struct{ Value string }

func (t *Text) Type() Type { return StringType }

// Inspect returns the text without quotes, because that is what a beginner
// expects `print "hi"` to put on screen.
func (t *Text) Inspect() string { return t.Value }

type Boolean struct{ Value bool }

func (b *Boolean) Type() Type { return BoolType }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type Nothing struct{}

func (n *Nothing) Type() Type      { return NothingType }
func (n *Nothing) Inspect() string { return "nothing" }

type Codeblock struct {
	Name       string
	Parameters []*ast.Identifier
	Body       *ast.Block
	Env        *Environment
}

func (c *Codeblock) Type() Type      { return CodeblockType }
func (c *Codeblock) Inspect() string { return "codeblock " + c.Name }

// ReturnValue carries a "give back" value up to the codeblock that should
// answer with it.
type ReturnValue struct{ Value Object }

func (r *ReturnValue) Type() Type      { return returnType }
func (r *ReturnValue) Inspect() string { return r.Value.Inspect() }

type Error struct {
	Message string
	Line    int
}

func (e *Error) Type() Type      { return errorType }
func (e *Error) Inspect() string { return fmt.Sprintf("line %d: %s", e.Line, e.Message) }

func IsError(obj Object) bool {
	_, ok := obj.(*Error)
	return ok
}

// Article picks "a" or "an" so error messages read like sentences.
func Article(t Type) string {
	switch t {
	case NumType:
		return "a num"
	case StringType:
		return "a string"
	case BoolType:
		return "a bool"
	case CodeblockType:
		return "a codeblock"
	}
	return "nothing"
}
