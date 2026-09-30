package parser

import (
	"strings"
	"testing"

	"github.com/abhinavtripathy/pharos/interpreter/ast"
)

func parse(t *testing.T, input string) *ast.Program {
	t.Helper()
	p := New(input)
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("unexpected parse errors for %q:\n  %s", input, strings.Join(errs, "\n  "))
	}
	return program
}

// parseErrors asserts that parsing fails and that the first error mentions
// want, so the messages stay aimed at a beginner rather than drifting into
// parser jargon.
func parseErrors(t *testing.T, input, want string) {
	t.Helper()
	p := New(input)
	p.ParseProgram()

	errs := p.Errors()
	if len(errs) == 0 {
		t.Fatalf("expected a parse error for %q, got none", input)
	}
	if !strings.Contains(errs[0], want) {
		t.Fatalf("error for %q was %q, want it to mention %q", input, errs[0], want)
	}
}

func onlyStatement(t *testing.T, input string) ast.Statement {
	t.Helper()
	program := parse(t, input)
	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement from %q, got %d", input, len(program.Statements))
	}
	return program.Statements[0]
}

func TestDeclare(t *testing.T) {
	cases := []struct{ input, want string }{
		{`num score = 7`, `num score = 7`},
		{`string name = "Ada"`, `string name = "Ada"`},
		{`bool done = false`, `bool done = false`},
		{`num total = 2 + 3 * 4`, `num total = (2 + (3 * 4))`},
	}
	for _, c := range cases {
		if got := onlyStatement(t, c.input).String(); got != c.want {
			t.Errorf("%q parsed as %q, want %q", c.input, got, c.want)
		}
	}
}

func TestAssign(t *testing.T) {
	if got := onlyStatement(t, `x = x + 1`).String(); got != `x = (x + 1)` {
		t.Errorf("got %q", got)
	}
}

func TestPrintAndAliases(t *testing.T) {
	for _, input := range []string{`print "hi"`, `out "hi"`, `output "hi"`} {
		stmt := onlyStatement(t, input)
		if _, ok := stmt.(*ast.Print); !ok {
			t.Fatalf("%q did not parse as a print statement, got %T", input, stmt)
		}
		if got := stmt.String(); got != `print "hi"` {
			t.Errorf("%q parsed as %q", input, got)
		}
	}
}

// Precedence is where a hand-written parser usually goes wrong, so pin the
// full shape with parentheses.
func TestOperatorPrecedence(t *testing.T) {
	cases := []struct{ input, want string }{
		{`print 1 + 2 + 3`, `print ((1 + 2) + 3)`},
		{`print 1 + 2 * 3`, `print (1 + (2 * 3))`},
		{`print 1 * 2 + 3`, `print ((1 * 2) + 3)`},
		{`print 8 / 4 / 2`, `print ((8 / 4) / 2)`},
		{`print (1 + 2) * 3`, `print ((1 + 2) * 3)`},
		{`print -a + b`, `print ((-a) + b)`},
		{`print 1 + 2 is greater than 3`, `print ((1 + 2) is greater than 3)`},
		{`print a is equal to b and c is less than d`, `print ((a is equal to b) and (c is less than d))`},
		{`print a and b or c`, `print ((a and b) or c)`},
		{`print not a and b`, `print ((not a) and b)`},
		{`print a is not equal to b`, `print (a is not equal to b)`},
	}
	for _, c := range cases {
		if got := onlyStatement(t, c.input).String(); got != c.want {
			t.Errorf("%q parsed as %q, want %q", c.input, got, c.want)
		}
	}
}

func TestIfChain(t *testing.T) {
	input := `if x is greater than 5 then
  print "big"
otherwise if x is equal to 5 then
  print "five"
otherwise
  print "small"
end if`

	stmt, ok := onlyStatement(t, input).(*ast.If)
	if !ok {
		t.Fatalf("did not parse as an if statement")
	}
	if len(stmt.Branches) != 2 {
		t.Fatalf("got %d branches, want 2", len(stmt.Branches))
	}
	if stmt.Else == nil {
		t.Fatal("expected an otherwise block")
	}
	if got := stmt.Branches[0].Condition.String(); got != `(x is greater than 5)` {
		t.Errorf("first condition was %q", got)
	}
}

func TestIfWithoutOtherwise(t *testing.T) {
	stmt, ok := onlyStatement(t, "if true then\n  print 1\nend if").(*ast.If)
	if !ok {
		t.Fatal("did not parse as an if statement")
	}
	if len(stmt.Branches) != 1 || stmt.Else != nil {
		t.Fatalf("got %d branches and else=%v", len(stmt.Branches), stmt.Else != nil)
	}
}

func TestNestedBlocks(t *testing.T) {
	input := `loop until done
  if x is less than 3 then
    x = x + 1
  otherwise
    done = true
  end if
end loop`

	stmt, ok := onlyStatement(t, input).(*ast.Loop)
	if !ok {
		t.Fatal("did not parse as a loop")
	}
	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("loop body has %d statements, want 1", len(stmt.Body.Statements))
	}
	if _, ok := stmt.Body.Statements[0].(*ast.If); !ok {
		t.Fatalf("loop body should hold an if, got %T", stmt.Body.Statements[0])
	}
}

func TestLoop(t *testing.T) {
	stmt, ok := onlyStatement(t, "loop until i is equal to 3\n  i = i + 1\nend loop").(*ast.Loop)
	if !ok {
		t.Fatal("did not parse as a loop")
	}
	if got := stmt.Condition.String(); got != `(i is equal to 3)` {
		t.Errorf("condition was %q", got)
	}
}

func TestCodeblockWithParameters(t *testing.T) {
	stmt, ok := onlyStatement(t, "codeblock add with a, b\n  give back a + b\nend codeblock").(*ast.Codeblock)
	if !ok {
		t.Fatal("did not parse as a codeblock")
	}
	if stmt.Name.Value != "add" {
		t.Errorf("name was %q", stmt.Name.Value)
	}
	if len(stmt.Parameters) != 2 {
		t.Fatalf("got %d parameters, want 2", len(stmt.Parameters))
	}
	if _, ok := stmt.Body.Statements[0].(*ast.GiveBack); !ok {
		t.Fatalf("body should hold a give back, got %T", stmt.Body.Statements[0])
	}
}

func TestCodeblockWithoutParameters(t *testing.T) {
	stmt, ok := onlyStatement(t, "codeblock hello\n  print \"hi\"\nend codeblock").(*ast.Codeblock)
	if !ok {
		t.Fatal("did not parse as a codeblock")
	}
	if len(stmt.Parameters) != 0 {
		t.Fatalf("got %d parameters, want 0", len(stmt.Parameters))
	}
}

func TestCallWithArguments(t *testing.T) {
	stmt := onlyStatement(t, `greet with "Ada", 3 + 4`)
	expr, ok := stmt.(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("got %T", stmt)
	}
	call, ok := expr.Expression.(*ast.Call)
	if !ok {
		t.Fatalf("got %T", expr.Expression)
	}
	if call.Function.Value != "greet" || len(call.Arguments) != 2 {
		t.Fatalf("parsed as %q", call.String())
	}
	if got := call.Arguments[1].String(); got != `(3 + 4)` {
		t.Errorf("second argument was %q", got)
	}
}

// A lone name on its own line is how you run a codeblock that takes nothing.
func TestBareNameIsACall(t *testing.T) {
	stmt := onlyStatement(t, `hello`)
	expr := stmt.(*ast.ExpressionStatement)
	call, ok := expr.Expression.(*ast.Call)
	if !ok {
		t.Fatalf("expected a call, got %T", expr.Expression)
	}
	if call.Function.Value != "hello" || len(call.Arguments) != 0 {
		t.Fatalf("parsed as %q", call.String())
	}
}

func TestCallUsedInsideExpression(t *testing.T) {
	if got := onlyStatement(t, `num x = double with 4`).String(); got != `num x = double with 4` {
		t.Errorf("got %q", got)
	}
	if got := onlyStatement(t, `print 1 + double with 4`).String(); got != `print (1 + double with 4)` {
		t.Errorf("got %q", got)
	}
}

func TestCommentsAndBlankLinesIgnored(t *testing.T) {
	input := `-- a leading note

num x = 1   -- trailing note

print x
`
	program := parse(t, input)
	if len(program.Statements) != 2 {
		t.Fatalf("got %d statements, want 2", len(program.Statements))
	}
}

func TestMultipleStatements(t *testing.T) {
	program := parse(t, "num a = 1\nnum b = 2\nprint a + b")
	if len(program.Statements) != 3 {
		t.Fatalf("got %d statements, want 3", len(program.Statements))
	}
}

func TestErrorsAreFriendly(t *testing.T) {
	cases := []struct{ input, want string }{
		{"num = 5", "a name"},
		{"num x 5", `"="`},
		{"if x then\n  print 1", `"end if"`},
		{"if x\n  print 1\nend if", `"then"`},
		{"loop until x\n  print 1", `"end loop"`},
		{"loop x\n print 1\nend loop", `"until"`},
		{"codeblock f\n  print 1", `"end codeblock"`},
		{"codeblock with a\n print 1\nend codeblock", "a name"},
		{"print", "needs something to print"},
		{`print "unclosed`, "never closed"},
		{"print 1 2", "unexpected"},
		{"print +", "I did not expect"},
	}
	for _, c := range cases {
		parseErrors(t, c.input, c.want)
	}
}

func TestErrorsIncludeLineNumbers(t *testing.T) {
	p := New("num a = 1\nnum b = 2\nnum = 3")
	p.ParseProgram()

	errs := p.Errors()
	if len(errs) == 0 {
		t.Fatal("expected an error")
	}
	if !strings.HasPrefix(errs[0], "line 3:") {
		t.Errorf("error should point at line 3, got %q", errs[0])
	}
}

// One broken line should not cascade into a pile of confusing follow-on
// errors.
func TestOneBadLineOneError(t *testing.T) {
	p := New("num x = 1\nnum = 2\nprint x")
	p.ParseProgram()

	if got := len(p.Errors()); got != 1 {
		t.Fatalf("got %d errors, want 1:\n  %s", got, strings.Join(p.Errors(), "\n  "))
	}
}

// String() should round-trip back into something the parser accepts again.
func TestStringRoundTrip(t *testing.T) {
	input := `codeblock classify with n
  if n is greater than 0 then
    give back "positive"
  otherwise if n is less than 0 then
    give back "negative"
  otherwise
    give back "zero"
  end if
end codeblock

print classify with -4`

	first := parse(t, input).String()
	second := parse(t, first).String()
	if first != second {
		t.Errorf("round trip changed the program:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}
