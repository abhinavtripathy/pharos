package lexer

import (
	"testing"

	"github.com/abhinavtripathy/pharos/interpreter/token"
)

type want struct {
	typ     token.Type
	literal string
}

func check(t *testing.T, input string, wants []want) {
	t.Helper()
	l := New(input)
	for i, w := range wants {
		tok := l.NextToken()
		if tok.Type != w.typ || tok.Literal != w.literal {
			t.Fatalf("token %d: got (%q, %q), want (%q, %q)",
				i, tok.Type, tok.Literal, w.typ, w.literal)
		}
	}
	if tok := l.NextToken(); tok.Type != token.EOF {
		t.Fatalf("expected EOF after %d tokens, got (%q, %q)", len(wants), tok.Type, tok.Literal)
	}
}

func TestDeclarations(t *testing.T) {
	check(t, `num score = 7`, []want{
		{token.NUM, "num"},
		{token.IDENT, "score"},
		{token.ASSIGN, "="},
		{token.NUMBER, "7"},
	})
}

func TestTextAndNumbers(t *testing.T) {
	check(t, `string s = "hi there"`, []want{
		{token.STRING, "string"},
		{token.IDENT, "s"},
		{token.ASSIGN, "="},
		{token.TEXT, "hi there"},
	})

	check(t, `3.5 + 2`, []want{
		{token.NUMBER, "3.5"},
		{token.PLUS, "+"},
		{token.NUMBER, "2"},
	})
}

func TestSingleQuotedText(t *testing.T) {
	check(t, `print 'hello'`, []want{
		{token.PRINT, "print"},
		{token.TEXT, "hello"},
	})
}

func TestTextEscapes(t *testing.T) {
	check(t, `"a\nb\t\"c\""`, []want{
		{token.TEXT, "a\nb\t\"c\""},
	})
}

func TestUnterminatedText(t *testing.T) {
	l := New("print \"oops\nprint 1")
	l.NextToken()
	if tok := l.NextToken(); tok.Type != token.ILLEGAL {
		t.Fatalf("expected ILLEGAL for unterminated text, got (%q, %q)", tok.Type, tok.Literal)
	}
}

// The multi-word operators are the heart of the language's readability, and
// the trickiest thing the lexer does.
func TestComparisonPhrases(t *testing.T) {
	check(t, `a is greater than b`, []want{
		{token.IDENT, "a"},
		{token.GT, "is greater than"},
		{token.IDENT, "b"},
	})

	check(t, `a is less than b`, []want{
		{token.IDENT, "a"},
		{token.LT, "is less than"},
		{token.IDENT, "b"},
	})

	check(t, `a is equal to b`, []want{
		{token.IDENT, "a"},
		{token.EQ, "is equal to"},
		{token.IDENT, "b"},
	})

	check(t, `a is not equal to b`, []want{
		{token.IDENT, "a"},
		{token.NEQ, "is not equal to"},
		{token.IDENT, "b"},
	})
}

func TestPhrasesWithIrregularSpacing(t *testing.T) {
	check(t, "a   is\tgreater  than   b", []want{
		{token.IDENT, "a"},
		{token.GT, "is\tgreater  than"},
		{token.IDENT, "b"},
	})
}

// "otherwise" and "otherwise if" overlap, so longest-match has to win.
func TestOtherwiseLongestMatch(t *testing.T) {
	check(t, "otherwise if x then", []want{
		{token.OTHERWISEIF, "otherwise if"},
		{token.IDENT, "x"},
		{token.THEN, "then"},
	})

	check(t, "otherwise\nprint 1", []want{
		{token.OTHERWISE, "otherwise"},
		{token.NEWLINE, "\n"},
		{token.PRINT, "print"},
		{token.NUMBER, "1"},
	})
}

// A phrase must never be stitched together across a line break.
func TestPhrasesDoNotCrossNewlines(t *testing.T) {
	check(t, "end\nif", []want{
		{token.IDENT, "end"},
		{token.NEWLINE, "\n"},
		{token.IF, "if"},
	})
}

func TestBlockKeywords(t *testing.T) {
	check(t, "end if", []want{{token.ENDIF, "end if"}})
	check(t, "end loop", []want{{token.ENDLOOP, "end loop"}})
	check(t, "end codeblock", []want{{token.ENDCODEBLOCK, "end codeblock"}})
	check(t, "give back 1", []want{{token.GIVEBACK, "give back"}, {token.NUMBER, "1"}})
}

func TestPrintAliases(t *testing.T) {
	check(t, "print 1", []want{{token.PRINT, "print"}, {token.NUMBER, "1"}})
	check(t, "out 1", []want{{token.PRINT, "out"}, {token.NUMBER, "1"}})
	check(t, "output 1", []want{{token.PRINT, "output"}, {token.NUMBER, "1"}})
}

func TestKeywordsAreCaseInsensitive(t *testing.T) {
	check(t, "IF x Is Greater Than 1 Then", []want{
		{token.IF, "IF"},
		{token.IDENT, "x"},
		{token.GT, "Is Greater Than"},
		{token.NUMBER, "1"},
		{token.THEN, "Then"},
	})
}

// An identifier that merely begins with a keyword's letters is still an
// identifier.
func TestIdentifiersNotSplitByKeywords(t *testing.T) {
	check(t, "iffy printer not_done", []want{
		{token.IDENT, "iffy"},
		{token.IDENT, "printer"},
		{token.IDENT, "not_done"},
	})
}

// "is" followed by something that is not a real operator stays an identifier
// rather than becoming a bogus comparison.
func TestDanglingIsIsIdentifier(t *testing.T) {
	check(t, "is flag", []want{
		{token.IDENT, "is"},
		{token.IDENT, "flag"},
	})
}

func TestComments(t *testing.T) {
	check(t, "-- just a note\nprint 1 -- trailing note", []want{
		{token.NEWLINE, "\n"},
		{token.PRINT, "print"},
		{token.NUMBER, "1"},
	})
}

// A comment directly after a word must not be swallowed into a phrase scan.
func TestCommentStopsPhraseLookahead(t *testing.T) {
	check(t, "otherwise -- a note\nprint 1", []want{
		{token.OTHERWISE, "otherwise"},
		{token.NEWLINE, "\n"},
		{token.PRINT, "print"},
		{token.NUMBER, "1"},
	})
}

func TestOperatorsAndPunctuation(t *testing.T) {
	check(t, "+ - * / = ( ) ,", []want{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.ASSIGN, "="},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.COMMA, ","},
	})
}

func TestLogicAndBooleans(t *testing.T) {
	check(t, "true and false or not true", []want{
		{token.TRUE, "true"},
		{token.AND, "and"},
		{token.FALSE, "false"},
		{token.OR, "or"},
		{token.NOT, "not"},
		{token.TRUE, "true"},
	})
}

func TestLineAndColumnTracking(t *testing.T) {
	l := New("print 1\nprint 2")
	tokens := l.Tokens()

	cases := []struct {
		index  int
		line   int
		column int
	}{
		{0, 1, 1},
		{1, 1, 7},
		{3, 2, 1},
		{4, 2, 7},
	}
	for _, c := range cases {
		got := tokens[c.index]
		if got.Line != c.line || got.Column != c.column {
			t.Errorf("token %d (%q): got line %d col %d, want line %d col %d",
				c.index, got.Literal, got.Line, got.Column, c.line, c.column)
		}
	}
}

func TestFullProgram(t *testing.T) {
	input := `num i = 0
loop until i is equal to 3
  print i
  i = i + 1
end loop`

	l := New(input)
	tokens := l.Tokens()
	for _, tok := range tokens {
		if tok.Type == token.ILLEGAL {
			t.Fatalf("unexpected ILLEGAL token: %q", tok.Literal)
		}
	}
	if got := tokens[len(tokens)-1].Type; got != token.EOF {
		t.Fatalf("expected trailing EOF, got %q", got)
	}
}
