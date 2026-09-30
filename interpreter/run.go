// Package interpreter is the front door to Pharos: it wires the lexer, parser
// and evaluator together so that callers only need to hand over source text.
package interpreter

import (
	"io"
	"strings"

	"github.com/abhinavtripathy/pharos/interpreter/evaluator"
	"github.com/abhinavtripathy/pharos/interpreter/object"
	"github.com/abhinavtripathy/pharos/interpreter/parser"
)

// Error collects everything that went wrong. Parsing reports as many problems
// as it can find at once; running stops at the first one.
type Error struct {
	Messages []string
}

func (e *Error) Error() string { return strings.Join(e.Messages, "\n") }

// Run executes a complete Pharos program, writing anything it prints to out.
func Run(source string, out io.Writer) error {
	return NewSession(out).Run(source)
}

// Session keeps variables and codeblocks alive across several pieces of
// source, which is what the interactive prompt needs.
type Session struct {
	env *object.Environment
	ev  *evaluator.Evaluator
}

func NewSession(out io.Writer) *Session {
	return &Session{env: object.NewEnvironment(), ev: evaluator.New(out)}
}

// Run evaluates source in the session, discarding the final value.
func (s *Session) Run(source string) error {
	_, err := s.Eval(source)
	return err
}

// Eval evaluates source and also hands back the value the last statement
// produced, so a prompt can echo it.
func (s *Session) Eval(source string) (object.Object, error) {
	p := parser.New(source)
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		return nil, &Error{Messages: errs}
	}

	result := s.ev.Eval(program, s.env)
	if failure, ok := result.(*object.Error); ok {
		return nil, &Error{Messages: []string{failure.Inspect()}}
	}
	return result, nil
}
