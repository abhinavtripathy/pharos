// Package parser builds an AST from Pharos source, collecting every error it
// finds rather than stopping at the first one.
package parser

import (
	"fmt"
	"strconv"

	"github.com/abhinavtripathy/pharos/interpreter/ast"
	"github.com/abhinavtripathy/pharos/interpreter/lexer"
	"github.com/abhinavtripathy/pharos/interpreter/token"
)

// Binding powers, loosest first.
const (
	lowest = iota + 1
	orPrec
	andPrec
	equals
	lessGreater
	sum
	product
	prefix
	call
)

var precedences = map[token.Type]int{
	token.OR:    orPrec,
	token.AND:   andPrec,
	token.EQ:    equals,
	token.NEQ:   equals,
	token.GT:    lessGreater,
	token.LT:    lessGreater,
	token.PLUS:  sum,
	token.MINUS: sum,
	token.STAR:  product,
	token.SLASH: product,
	token.WITH:  call,
}

type Parser struct {
	l    *lexer.Lexer
	cur  token.Token
	peek token.Token

	errors []string
}

func New(input string) *Parser {
	p := &Parser{l: lexer.New(input)}
	p.nextToken()
	p.nextToken()
	return p
}

// Errors returns every problem found, each prefixed with its line number.
func (p *Parser) Errors() []string { return p.errors }

func (p *Parser) nextToken() {
	p.cur = p.peek
	p.peek = p.l.NextToken()
}

func (p *Parser) curIs(types ...token.Type) bool {
	for _, t := range types {
		if p.cur.Type == t {
			return true
		}
	}
	return false
}

func (p *Parser) peekIs(t token.Type) bool { return p.peek.Type == t }

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekIs(t) {
		p.nextToken()
		return true
	}
	p.errorf(p.peek, "expected %s here, but found %s", describe(t), found(p.peek))
	return false
}

func (p *Parser) errorf(tok token.Token, format string, args ...any) {
	p.errors = append(p.errors, fmt.Sprintf("line %d: %s", tok.Line, fmt.Sprintf(format, args...)))
}

func (p *Parser) skipNewlines() {
	for p.curIs(token.NEWLINE) {
		p.nextToken()
	}
}

// ParseProgram consumes the whole input. Statements are separated by line
// breaks; there are no semicolons in Pharos.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	p.skipNewlines()

	for !p.curIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.endStatement(stmt != nil)
	}
	return program
}

// endStatement moves past the end of a statement. When the statement itself
// parsed cleanly, leftover tokens on the line are worth reporting; when it did
// not, we skip quietly to avoid piling a second error onto the first.
func (p *Parser) endStatement(parsed bool) {
	if !p.curIs(token.NEWLINE, token.EOF) {
		if p.peekIs(token.NEWLINE) || p.peekIs(token.EOF) {
			p.nextToken()
		} else {
			if parsed {
				p.errorf(p.cur, "unexpected %s at the end of this line", found(p.cur))
			}
			for !p.curIs(token.NEWLINE, token.EOF) {
				p.nextToken()
			}
		}
	}
	p.skipNewlines()
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.cur.Type {
	case token.NUM, token.STRING, token.BOOL:
		return p.parseDeclare()
	case token.PRINT:
		return p.parsePrint()
	case token.IF:
		return p.parseIf()
	case token.LOOP:
		return p.parseLoop()
	case token.CODEBLOCK:
		return p.parseCodeblock()
	case token.GIVEBACK:
		return p.parseGiveBack()
	case token.IDENT:
		if p.peekIs(token.ASSIGN) {
			return p.parseAssign()
		}
	case token.ILLEGAL:
		p.errorf(p.cur, "%s", p.cur.Literal)
		return nil
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseDeclare() ast.Statement {
	stmt := &ast.Declare{Tok: p.cur, DeclaredType: string(p.cur.Type)}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Tok: p.cur, Value: p.cur.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(lowest)
	if stmt.Value == nil {
		return nil
	}
	return stmt
}

func (p *Parser) parseAssign() ast.Statement {
	stmt := &ast.Assign{Tok: p.cur, Name: &ast.Identifier{Tok: p.cur, Value: p.cur.Literal}}

	p.nextToken() // the "="
	p.nextToken()

	stmt.Value = p.parseExpression(lowest)
	if stmt.Value == nil {
		return nil
	}
	return stmt
}

func (p *Parser) parsePrint() ast.Statement {
	stmt := &ast.Print{Tok: p.cur}

	if p.peekIs(token.NEWLINE) || p.peekIs(token.EOF) {
		p.errorf(p.cur, "%q needs something to print after it", p.cur.Literal)
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(lowest)
	if stmt.Value == nil {
		return nil
	}
	return stmt
}

func (p *Parser) parseIf() ast.Statement {
	stmt := &ast.If{Tok: p.cur}

	branch, ok := p.parseBranch()
	if !ok {
		return nil
	}
	stmt.Branches = append(stmt.Branches, branch)

	for p.curIs(token.OTHERWISEIF) {
		branch, ok := p.parseBranch()
		if !ok {
			return nil
		}
		stmt.Branches = append(stmt.Branches, branch)
	}

	if p.curIs(token.OTHERWISE) {
		if !p.expectPeek(token.NEWLINE) {
			return nil
		}
		stmt.Else = p.parseBlock(token.ENDIF)
	}

	if !p.curIs(token.ENDIF) {
		p.errorf(p.cur, "this %q is missing its %q", stmt.Tok.Literal, "end if")
		return nil
	}
	return stmt
}

// parseBranch reads "<cond> then <newline> <body>", starting from the "if" or
// "otherwise if" keyword and stopping on whichever keyword closes the body.
func (p *Parser) parseBranch() (ast.Branch, bool) {
	opener := p.cur
	p.nextToken()

	cond := p.parseExpression(lowest)
	if cond == nil {
		return ast.Branch{}, false
	}

	if !p.peekIs(token.THEN) {
		p.errorf(p.peek, "the condition after %q must be followed by %q, but I found %s",
			opener.Literal, "then", found(p.peek))
		return ast.Branch{}, false
	}
	p.nextToken()

	if !p.expectPeek(token.NEWLINE) {
		return ast.Branch{}, false
	}

	body := p.parseBlock(token.OTHERWISEIF, token.OTHERWISE, token.ENDIF)
	return ast.Branch{Condition: cond, Body: body}, true
}

func (p *Parser) parseLoop() ast.Statement {
	stmt := &ast.Loop{Tok: p.cur}

	if !p.expectPeek(token.UNTIL) {
		return nil
	}
	p.nextToken()

	stmt.Condition = p.parseExpression(lowest)
	if stmt.Condition == nil {
		return nil
	}

	if !p.expectPeek(token.NEWLINE) {
		return nil
	}

	stmt.Body = p.parseBlock(token.ENDLOOP)
	if !p.curIs(token.ENDLOOP) {
		p.errorf(stmt.Tok, "this %q is missing its %q", stmt.Tok.Literal, "end loop")
		return nil
	}
	return stmt
}

func (p *Parser) parseCodeblock() ast.Statement {
	stmt := &ast.Codeblock{Tok: p.cur}

	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Tok: p.cur, Value: p.cur.Literal}

	if p.peekIs(token.WITH) {
		p.nextToken()
		for {
			if !p.expectPeek(token.IDENT) {
				return nil
			}
			stmt.Parameters = append(stmt.Parameters, &ast.Identifier{Tok: p.cur, Value: p.cur.Literal})
			if !p.peekIs(token.COMMA) {
				break
			}
			p.nextToken()
		}
	}

	if !p.expectPeek(token.NEWLINE) {
		return nil
	}

	stmt.Body = p.parseBlock(token.ENDCODEBLOCK)
	if !p.curIs(token.ENDCODEBLOCK) {
		p.errorf(stmt.Tok, "this %q is missing its %q", stmt.Tok.Literal, "end codeblock")
		return nil
	}
	return stmt
}

func (p *Parser) parseGiveBack() ast.Statement {
	stmt := &ast.GiveBack{Tok: p.cur}
	if p.peekIs(token.NEWLINE) || p.peekIs(token.EOF) {
		return stmt
	}
	p.nextToken()

	stmt.Value = p.parseExpression(lowest)
	if stmt.Value == nil {
		return nil
	}
	return stmt
}

// parseBlock reads statements until it meets one of the given closing
// keywords, leaving that keyword as the current token.
func (p *Parser) parseBlock(terminators ...token.Type) *ast.Block {
	block := &ast.Block{Tok: p.cur}
	p.nextToken()
	p.skipNewlines()

	for !p.curIs(terminators...) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.endStatement(stmt != nil)
	}
	return block
}

// parseExpressionStatement also handles a bare name on a line of its own,
// which is how a codeblock that takes no values is run.
func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Tok: p.cur}

	expr := p.parseExpression(lowest)
	if expr == nil {
		return nil
	}

	if ident, ok := expr.(*ast.Identifier); ok {
		expr = &ast.Call{Tok: ident.Tok, Function: ident}
	}
	stmt.Expression = expr
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	left := p.parsePrefixExpression()
	if left == nil {
		return nil
	}

	for precedence < p.peekPrecedence() {
		p.nextToken()
		left = p.parseInfixExpression(left)
		if left == nil {
			return nil
		}
	}
	return left
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peek.Type]; ok {
		return prec
	}
	return lowest
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.cur.Type]; ok {
		return prec
	}
	return lowest
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	switch p.cur.Type {
	case token.IDENT:
		return &ast.Identifier{Tok: p.cur, Value: p.cur.Literal}

	case token.NUMBER:
		value, err := strconv.ParseFloat(p.cur.Literal, 64)
		if err != nil {
			p.errorf(p.cur, "%q is not a number I can understand", p.cur.Literal)
			return nil
		}
		return &ast.NumberLiteral{Tok: p.cur, Value: value}

	case token.TEXT:
		return &ast.TextLiteral{Tok: p.cur, Value: p.cur.Literal}

	case token.TRUE, token.FALSE:
		return &ast.BooleanLiteral{Tok: p.cur, Value: p.curIs(token.TRUE)}

	case token.NOT, token.MINUS:
		expr := &ast.Prefix{Tok: p.cur, Operator: string(p.cur.Type)}
		p.nextToken()
		expr.Right = p.parseExpression(prefix)
		if expr.Right == nil {
			return nil
		}
		return expr

	case token.LPAREN:
		p.nextToken()
		inner := p.parseExpression(lowest)
		if inner == nil {
			return nil
		}
		if !p.expectPeek(token.RPAREN) {
			return nil
		}
		return inner

	case token.ILLEGAL:
		p.errorf(p.cur, "%s", p.cur.Literal)
		return nil
	}

	p.errorf(p.cur, "I did not expect %s here", found(p.cur))
	return nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	if p.curIs(token.WITH) {
		return p.parseCall(left)
	}

	expr := &ast.Infix{Tok: p.cur, Left: left, Operator: string(p.cur.Type)}
	precedence := p.curPrecedence()
	p.nextToken()

	expr.Right = p.parseExpression(precedence)
	if expr.Right == nil {
		return nil
	}
	return expr
}

func (p *Parser) parseCall(left ast.Expression) ast.Expression {
	ident, ok := left.(*ast.Identifier)
	if !ok {
		p.errorf(p.cur, "only the name of a codeblock can come before %q", "with")
		return nil
	}

	expr := &ast.Call{Tok: p.cur, Function: ident}
	p.nextToken()

	for {
		arg := p.parseExpression(lowest)
		if arg == nil {
			return nil
		}
		expr.Arguments = append(expr.Arguments, arg)

		if !p.peekIs(token.COMMA) {
			return expr
		}
		p.nextToken() // the comma
		p.nextToken()
	}
}

// describe names a token type the way the guide names it, for error messages
// aimed at someone who has never seen a parser error before.
func describe(t token.Type) string {
	switch t {
	case token.IDENT:
		return "a name"
	case token.NUMBER:
		return "a number"
	case token.TEXT:
		return "a piece of text in quotes"
	case token.NEWLINE:
		return "a new line"
	case token.EOF:
		return "the end of the program"
	}
	return strconv.Quote(string(t))
}

// found describes the token actually sitting in front of the parser.
func found(tok token.Token) string {
	switch tok.Type {
	case token.EOF:
		return "the end of the program"
	case token.NEWLINE:
		return "the end of the line"
	case token.TEXT:
		return "the text " + strconv.Quote(tok.Literal)
	}
	return strconv.Quote(tok.Literal)
}
