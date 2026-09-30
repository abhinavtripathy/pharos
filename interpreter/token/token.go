// Package token defines the lexical vocabulary of the Pharos language.
package token

import "strings"

type Type string

// Token is a single lexical unit along with where it was found, so that
// later stages can point a beginner at the exact line that went wrong.
type Token struct {
	Type    Type
	Literal string
	Line    int
	Column  int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	NEWLINE = "NEWLINE"

	IDENT  = "IDENT"
	NUMBER = "NUMBER"
	TEXT   = "TEXT"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	STAR   = "*"
	SLASH  = "/"

	LPAREN = "("
	RPAREN = ")"
	COMMA  = ","

	// Declaration types.
	NUM    = "num"
	STRING = "string"
	BOOL   = "bool"
	TRUE   = "true"
	FALSE  = "false"

	PRINT = "print"

	IF          = "if"
	THEN        = "then"
	OTHERWISE   = "otherwise"
	OTHERWISEIF = "otherwise if"
	ENDIF       = "end if"

	LOOP    = "loop"
	UNTIL   = "until"
	ENDLOOP = "end loop"

	CODEBLOCK    = "codeblock"
	ENDCODEBLOCK = "end codeblock"
	WITH         = "with"
	GIVEBACK     = "give back"

	AND = "and"
	OR  = "or"
	NOT = "not"

	GT  = "is greater than"
	LT  = "is less than"
	EQ  = "is equal to"
	NEQ = "is not equal to"
)

// phrases maps every keyword to its token type. Multi-word entries are what
// make Pharos read like English; the lexer resolves them by longest match so
// the parser only ever sees a single token for "is greater than".
var phrases = map[string]Type{
	"num":    NUM,
	"string": STRING,
	"bool":   BOOL,
	"true":   TRUE,
	"false":  FALSE,

	"print":  PRINT,
	"out":    PRINT,
	"output": PRINT,

	"if":           IF,
	"then":         THEN,
	"otherwise":    OTHERWISE,
	"otherwise if": OTHERWISEIF,
	"end if":       ENDIF,

	"loop":     LOOP,
	"until":    UNTIL,
	"end loop": ENDLOOP,

	"codeblock":     CODEBLOCK,
	"end codeblock": ENDCODEBLOCK,
	"with":          WITH,
	"give back":     GIVEBACK,

	"and": AND,
	"or":  OR,
	"not": NOT,

	"is greater than": GT,
	"is less than":    LT,
	"is equal to":     EQ,
	"is not equal to": NEQ,
}

// MaxPhraseWords is the longest keyword in words ("is not equal to"), and so
// the furthest the lexer ever needs to look ahead.
const MaxPhraseWords = 4

// LookupPhrase resolves a run of words to a keyword token type. Keywords are
// matched without regard to case, which is friendlier for newcomers; variable
// names stay case-sensitive.
func LookupPhrase(words []string) (Type, bool) {
	t, ok := phrases[strings.ToLower(strings.Join(words, " "))]
	return t, ok
}
