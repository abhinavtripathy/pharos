package token

import (
	"strings"
	"testing"
)

// The lexer only ever looks MaxPhraseWords ahead. A longer keyword added to
// the table would therefore never match, and would do so silently, so the two
// are pinned together here.
func TestNoPhraseIsLongerThanTheLookahead(t *testing.T) {
	for phrase := range phrases {
		if words := len(strings.Fields(phrase)); words > MaxPhraseWords {
			t.Errorf("keyword %q is %d words, but the lexer only looks %d ahead",
				phrase, words, MaxPhraseWords)
		}
	}
}

// Every keyword in the table must be lower case, since lookups are made after
// lowercasing.
func TestPhraseTableIsLowerCase(t *testing.T) {
	for phrase := range phrases {
		if phrase != strings.ToLower(phrase) {
			t.Errorf("keyword %q should be written in lower case in the table", phrase)
		}
	}
}

func TestLookupPhrase(t *testing.T) {
	cases := []struct {
		words []string
		want  Type
		found bool
	}{
		{[]string{"num"}, NUM, true},
		{[]string{"IF"}, IF, true},
		{[]string{"Is", "Greater", "Than"}, GT, true},
		{[]string{"is", "not", "equal", "to"}, NEQ, true},
		{[]string{"otherwise", "if"}, OTHERWISEIF, true},
		{[]string{"out"}, PRINT, true},
		{[]string{"output"}, PRINT, true},
		{[]string{"end"}, "", false},
		{[]string{"is"}, "", false},
		{[]string{"score"}, "", false},
	}

	for _, c := range cases {
		got, found := LookupPhrase(c.words)
		if found != c.found || got != c.want {
			t.Errorf("LookupPhrase(%v) = (%q, %v), want (%q, %v)",
				c.words, got, found, c.want, c.found)
		}
	}
}
