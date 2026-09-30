// Package ast describes the shape of a parsed Pharos program.
package ast

import (
	"strconv"
	"strings"

	"github.com/abhinavtripathy/pharos/interpreter/token"
)

type Node interface {
	Token() token.Token
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Token() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].Token()
	}
	return token.Token{}
}

func (p *Program) String() string {
	var sb strings.Builder
	for _, s := range p.Statements {
		sb.WriteString(s.String())
		sb.WriteString("\n")
	}
	return sb.String()
}

type Block struct {
	Tok        token.Token
	Statements []Statement
}

func (b *Block) Token() token.Token { return b.Tok }
func (b *Block) String() string {
	var sb strings.Builder
	for _, s := range b.Statements {
		sb.WriteString(indent(s.String()))
		sb.WriteString("\n")
	}
	return sb.String()
}

// Declare is "num x = 5". DeclaredType is the type keyword the user wrote,
// which the evaluator holds them to on every later assignment.
type Declare struct {
	Tok          token.Token
	DeclaredType string
	Name         *Identifier
	Value        Expression
}

func (d *Declare) Token() token.Token { return d.Tok }
func (d *Declare) statementNode()     {}
func (d *Declare) String() string {
	return d.DeclaredType + " " + d.Name.String() + " = " + d.Value.String()
}

// Assign is "x = 5" for a variable that already exists.
type Assign struct {
	Tok   token.Token
	Name  *Identifier
	Value Expression
}

func (a *Assign) Token() token.Token { return a.Tok }
func (a *Assign) statementNode()     {}
func (a *Assign) String() string     { return a.Name.String() + " = " + a.Value.String() }

type Print struct {
	Tok   token.Token
	Value Expression
}

func (p *Print) Token() token.Token { return p.Tok }
func (p *Print) statementNode()     {}
func (p *Print) String() string     { return "print " + p.Value.String() }

// Branch is one "if"/"otherwise if" arm of a conditional.
type Branch struct {
	Condition Expression
	Body      *Block
}

type If struct {
	Tok      token.Token
	Branches []Branch
	Else     *Block
}

func (i *If) Token() token.Token { return i.Tok }
func (i *If) statementNode()     {}
func (i *If) String() string {
	var sb strings.Builder
	for n, br := range i.Branches {
		if n == 0 {
			sb.WriteString("if ")
		} else {
			sb.WriteString("otherwise if ")
		}
		sb.WriteString(br.Condition.String())
		sb.WriteString(" then\n")
		sb.WriteString(br.Body.String())
	}
	if i.Else != nil {
		sb.WriteString("otherwise\n")
		sb.WriteString(i.Else.String())
	}
	sb.WriteString("end if")
	return sb.String()
}

// Loop is "loop until <condition>": it repeats while the condition is false.
type Loop struct {
	Tok       token.Token
	Condition Expression
	Body      *Block
}

func (l *Loop) Token() token.Token { return l.Tok }
func (l *Loop) statementNode()     {}
func (l *Loop) String() string {
	return "loop until " + l.Condition.String() + "\n" + l.Body.String() + "end loop"
}

type Codeblock struct {
	Tok        token.Token
	Name       *Identifier
	Parameters []*Identifier
	Body       *Block
}

func (c *Codeblock) Token() token.Token { return c.Tok }
func (c *Codeblock) statementNode()     {}
func (c *Codeblock) String() string {
	var sb strings.Builder
	sb.WriteString("codeblock " + c.Name.String())
	if len(c.Parameters) > 0 {
		names := make([]string, len(c.Parameters))
		for i, p := range c.Parameters {
			names[i] = p.String()
		}
		sb.WriteString(" with " + strings.Join(names, ", "))
	}
	sb.WriteString("\n")
	sb.WriteString(c.Body.String())
	sb.WriteString("end codeblock")
	return sb.String()
}

type GiveBack struct {
	Tok   token.Token
	Value Expression
}

func (g *GiveBack) Token() token.Token { return g.Tok }
func (g *GiveBack) statementNode()     {}
func (g *GiveBack) String() string {
	if g.Value == nil {
		return "give back"
	}
	return "give back " + g.Value.String()
}

type ExpressionStatement struct {
	Tok        token.Token
	Expression Expression
}

func (e *ExpressionStatement) Token() token.Token { return e.Tok }
func (e *ExpressionStatement) statementNode()     {}
func (e *ExpressionStatement) String() string     { return e.Expression.String() }

type Identifier struct {
	Tok   token.Token
	Value string
}

func (i *Identifier) Token() token.Token { return i.Tok }
func (i *Identifier) expressionNode()    {}
func (i *Identifier) String() string     { return i.Value }

type NumberLiteral struct {
	Tok   token.Token
	Value float64
}

func (n *NumberLiteral) Token() token.Token { return n.Tok }
func (n *NumberLiteral) expressionNode()    {}
func (n *NumberLiteral) String() string     { return FormatNumber(n.Value) }

type TextLiteral struct {
	Tok   token.Token
	Value string
}

func (t *TextLiteral) Token() token.Token { return t.Tok }
func (t *TextLiteral) expressionNode()    {}
func (t *TextLiteral) String() string     { return strconv.Quote(t.Value) }

type BooleanLiteral struct {
	Tok   token.Token
	Value bool
}

func (b *BooleanLiteral) Token() token.Token { return b.Tok }
func (b *BooleanLiteral) expressionNode()    {}
func (b *BooleanLiteral) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type Prefix struct {
	Tok      token.Token
	Operator string
	Right    Expression
}

func (p *Prefix) Token() token.Token { return p.Tok }
func (p *Prefix) expressionNode()    {}
func (p *Prefix) String() string {
	if p.Operator == "-" {
		return "(-" + p.Right.String() + ")"
	}
	return "(" + p.Operator + " " + p.Right.String() + ")"
}

type Infix struct {
	Tok      token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *Infix) Token() token.Token { return i.Tok }
func (i *Infix) expressionNode()    {}
func (i *Infix) String() string {
	return "(" + i.Left.String() + " " + i.Operator + " " + i.Right.String() + ")"
}

// Call is "greet with 1, 2", or a bare "greet" used as a statement.
type Call struct {
	Tok       token.Token
	Function  *Identifier
	Arguments []Expression
}

func (c *Call) Token() token.Token { return c.Tok }
func (c *Call) expressionNode()    {}
func (c *Call) String() string {
	if len(c.Arguments) == 0 {
		return c.Function.String()
	}
	args := make([]string, len(c.Arguments))
	for i, a := range c.Arguments {
		args[i] = a.String()
	}
	return c.Function.String() + " with " + strings.Join(args, ", ")
}

// FormatNumber renders a number the way a beginner expects: 4 rather than 4.0,
// but 4.5 still shows its fraction.
func FormatNumber(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func indent(s string) string {
	s = strings.TrimRight(s, "\n")
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = "  " + l
	}
	return strings.Join(lines, "\n")
}
