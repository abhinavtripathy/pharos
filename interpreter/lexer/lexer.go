// Package lexer turns Pharos source text into a stream of tokens.
package lexer

import (
	"strings"

	"github.com/abhinavtripathy/pharos/interpreter/token"
)

type Lexer struct {
	input     string
	pos       int
	line      int
	lineStart int
}

func New(input string) *Lexer {
	return &Lexer{input: input, line: 1}
}

// NextToken returns the next token, or an EOF token once the input is spent.
func (l *Lexer) NextToken() token.Token {
	l.skipSpaceAndComments()

	if l.pos >= len(l.input) {
		return l.make(token.EOF, "")
	}

	ch := l.input[l.pos]

	switch {
	case ch == '\n':
		tok := l.make(token.NEWLINE, "\n")
		l.pos++
		l.line++
		l.lineStart = l.pos
		return tok
	case isDigit(ch):
		return l.readNumber()
	case ch == '"' || ch == '\'':
		return l.readText()
	case isWordStart(ch):
		return l.readWords()
	}

	tok := l.make(singleCharType(ch), string(ch))
	l.pos++
	return tok
}

// Tokens drains the lexer. Handy for tests and for the REPL, which counts
// block openers and closers to decide whether an entry is complete.
func (l *Lexer) Tokens() []token.Token {
	var out []token.Token
	for {
		tok := l.NextToken()
		out = append(out, tok)
		if tok.Type == token.EOF {
			return out
		}
	}
}

func (l *Lexer) make(t token.Type, literal string) token.Token {
	return token.Token{Type: t, Literal: literal, Line: l.line, Column: l.pos - l.lineStart + 1}
}

func (l *Lexer) skipSpaceAndComments() {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\r':
			l.pos++
		case '-':
			if !l.commentAt(l.pos) {
				return
			}
			for l.pos < len(l.input) && l.input[l.pos] != '\n' {
				l.pos++
			}
		default:
			return
		}
	}
}

func (l *Lexer) commentAt(i int) bool {
	return i+1 < len(l.input) && l.input[i] == '-' && l.input[i+1] == '-'
}

func (l *Lexer) readNumber() token.Token {
	start := l.pos
	tok := l.make(token.NUMBER, "")
	for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
		l.pos++
	}
	if l.pos+1 < len(l.input) && l.input[l.pos] == '.' && isDigit(l.input[l.pos+1]) {
		l.pos++
		for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
			l.pos++
		}
	}
	tok.Literal = l.input[start:l.pos]
	return tok
}

// readText reads a quoted string. Both quote styles are accepted because both
// were highlighted in the original editor.
func (l *Lexer) readText() token.Token {
	tok := l.make(token.TEXT, "")
	quote := l.input[l.pos]
	l.pos++

	var sb strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		switch {
		case ch == quote:
			l.pos++
			tok.Literal = sb.String()
			return tok
		case ch == '\n':
			tok.Type = token.ILLEGAL
			tok.Literal = "a piece of text was never closed with " + string(quote)
			return tok
		case ch == '\\' && l.pos+1 < len(l.input):
			l.pos++
			sb.WriteByte(unescape(l.input[l.pos]))
			l.pos++
		default:
			sb.WriteByte(ch)
			l.pos++
		}
	}

	tok.Type = token.ILLEGAL
	tok.Literal = "a piece of text was never closed with " + string(quote)
	return tok
}

func unescape(ch byte) byte {
	switch ch {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	default:
		return ch
	}
}

// readWords consumes one word, then looks ahead over the following words to
// find the longest keyword phrase that matches. Anything that matches nothing
// is a plain identifier.
func (l *Lexer) readWords() token.Token {
	tok := l.make(token.IDENT, "")

	words := make([]string, 0, token.MaxPhraseWords)
	ends := make([]int, 0, token.MaxPhraseWords)

	start := l.pos
	cursor := l.pos
	for len(words) < token.MaxPhraseWords {
		word, end := l.wordAt(cursor)
		if word == "" {
			break
		}
		words = append(words, word)
		ends = append(ends, end)

		cursor = end
		for cursor < len(l.input) && (l.input[cursor] == ' ' || l.input[cursor] == '\t') {
			cursor++
		}
		if l.commentAt(cursor) {
			break
		}
	}

	for n := len(words); n >= 1; n-- {
		if t, ok := token.LookupPhrase(words[:n]); ok {
			l.pos = ends[n-1]
			tok.Type = t
			tok.Literal = l.input[start:l.pos]
			return tok
		}
	}

	l.pos = ends[0]
	tok.Literal = l.input[start:l.pos]
	return tok
}

func (l *Lexer) wordAt(i int) (string, int) {
	if i >= len(l.input) || !isWordStart(l.input[i]) {
		return "", i
	}
	j := i
	for j < len(l.input) && isWordPart(l.input[j]) {
		j++
	}
	return l.input[i:j], j
}

func singleCharType(ch byte) token.Type {
	switch ch {
	case '=':
		return token.ASSIGN
	case '+':
		return token.PLUS
	case '-':
		return token.MINUS
	case '*':
		return token.STAR
	case '/':
		return token.SLASH
	case '(':
		return token.LPAREN
	case ')':
		return token.RPAREN
	case ',':
		return token.COMMA
	}
	return token.ILLEGAL
}

func isDigit(ch byte) bool     { return ch >= '0' && ch <= '9' }
func isWordStart(ch byte) bool { return ch == '_' || (ch|0x20 >= 'a' && ch|0x20 <= 'z') }
func isWordPart(ch byte) bool  { return isWordStart(ch) || isDigit(ch) }
