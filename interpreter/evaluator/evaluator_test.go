package evaluator_test

import (
	"strings"
	"testing"

	"github.com/abhinavtripathy/pharos/interpreter"
)

// run executes source and returns everything it printed.
func run(t *testing.T, source string) string {
	t.Helper()
	var out strings.Builder
	if err := interpreter.Run(source, &out); err != nil {
		t.Fatalf("unexpected error running:\n%s\n\ngot: %v", source, err)
	}
	return out.String()
}

// expectPrints asserts the printed output, with a trailing newline implied.
func expectPrints(t *testing.T, source, want string) {
	t.Helper()
	got := run(t, source)
	if got != want+"\n" {
		t.Errorf("%s\n  printed %q, want %q", source, got, want+"\n")
	}
}

// expectFails asserts the program fails with a message mentioning want.
func expectFails(t *testing.T, source, want string) {
	t.Helper()
	var out strings.Builder
	err := interpreter.Run(source, &out)
	if err == nil {
		t.Fatalf("expected an error from:\n%s", source)
	}
	if !strings.Contains(err.Error(), want) {
		t.Errorf("%s\n  failed with %q, want it to mention %q", source, err.Error(), want)
	}
}

func TestArithmetic(t *testing.T) {
	cases := []struct{ source, want string }{
		{`print 1 + 2`, "3"},
		{`print 7 - 9`, "-2"},
		{`print 3 * 4`, "12"},
		{`print 10 / 4`, "2.5"},
		{`print 2 + 3 * 4`, "14"},
		{`print (2 + 3) * 4`, "20"},
		{`print -5 + 2`, "-3"},
		{`print 1.5 + 1.5`, "3"},
	}
	for _, c := range cases {
		expectPrints(t, c.source, c.want)
	}
}

// Whole numbers should not grow a decimal point just because they are stored
// as floats.
func TestNumbersPrintCleanly(t *testing.T) {
	expectPrints(t, `print 4`, "4")
	expectPrints(t, `print 8 / 2`, "4")
	expectPrints(t, `print 7 / 2`, "3.5")
}

func TestDivideByZero(t *testing.T) {
	expectFails(t, `print 1 / 0`, "cannot divide by zero")
}

func TestTextJoining(t *testing.T) {
	expectPrints(t, `print "Hello " + "world"`, "Hello world")
	expectPrints(t, `print "score: " + 10`, "score: 10")
	expectPrints(t, `print 10 + " points"`, "10 points")
	expectPrints(t, `print "done: " + true`, "done: true")
}

func TestTextPrintsWithoutQuotes(t *testing.T) {
	expectPrints(t, `print "hi"`, "hi")
}

func TestDeclarationsAndAssignment(t *testing.T) {
	expectPrints(t, "num x = 5\nprint x", "5")
	expectPrints(t, "num x = 5\nx = x + 1\nprint x", "6")
	expectPrints(t, "string s = \"a\"\ns = s + \"b\"\nprint s", "ab")
	expectPrints(t, "bool b = false\nb = not b\nprint b", "true")
}

// The declared type is the whole point of writing num/string/bool, so it has
// to actually be enforced.
func TestDeclaredTypesAreEnforced(t *testing.T) {
	expectFails(t, `num score = "seven"`, "cannot hold a string")
	expectFails(t, `string name = 4`, "cannot hold a num")
	expectFails(t, `bool done = 1`, "cannot hold a num")
	expectFails(t, "num score = 1\nscore = \"two\"", "cannot hold a string")
}

func TestUndeclaredNames(t *testing.T) {
	expectFails(t, `print mystery`, "have not heard of")
	expectFails(t, `x = 5`, "has not been declared yet")
}

func TestRedeclaration(t *testing.T) {
	expectFails(t, "num x = 1\nnum x = 2", "already been declared")
}

func TestComparisons(t *testing.T) {
	cases := []struct{ source, want string }{
		{`print 5 is greater than 3`, "true"},
		{`print 3 is greater than 5`, "false"},
		{`print 3 is less than 5`, "true"},
		{`print 4 is equal to 4`, "true"},
		{`print 4 is not equal to 4`, "false"},
		{`print "a" is equal to "a"`, "true"},
		{`print "a" is not equal to "b"`, "true"},
		{`print true is equal to true`, "true"},
	}
	for _, c := range cases {
		expectPrints(t, c.source, c.want)
	}
}

func TestComparisonTypeErrors(t *testing.T) {
	expectFails(t, `print "a" is greater than 1`, "only works with nums")
	expectFails(t, `print 1 is equal to "a"`, "cannot compare")
}

func TestLogic(t *testing.T) {
	cases := []struct{ source, want string }{
		{`print true and true`, "true"},
		{`print true and false`, "false"},
		{`print false or true`, "true"},
		{`print false or false`, "false"},
		{`print not true`, "false"},
		{`print 1 is less than 2 and 2 is less than 3`, "true"},
	}
	for _, c := range cases {
		expectPrints(t, c.source, c.want)
	}
}

// "and" must not evaluate its right side once the left is false, or a guard
// like `count is not equal to 0 and total / count is greater than 1` would
// still divide by zero.
func TestLogicStopsEarly(t *testing.T) {
	expectPrints(t, `print false and 1 / 0 is equal to 1`, "false")
	expectPrints(t, `print true or 1 / 0 is equal to 1`, "true")
}

func TestLogicTypeErrors(t *testing.T) {
	expectFails(t, `print 1 and true`, "needs true or false on its left")
	expectFails(t, `print true and 1`, "needs true or false on its right")
	expectFails(t, `print not 1`, "needs true or false")
}

func TestIfChain(t *testing.T) {
	program := `num score = %s
if score is greater than 5 then
  print "big"
otherwise if score is equal to 5 then
  print "five"
otherwise
  print "small"
end if`

	expectPrints(t, strings.Replace(program, "%s", "9", 1), "big")
	expectPrints(t, strings.Replace(program, "%s", "5", 1), "five")
	expectPrints(t, strings.Replace(program, "%s", "1", 1), "small")
}

func TestIfWithNoMatchingBranch(t *testing.T) {
	if got := run(t, "if false then\n  print \"no\"\nend if"); got != "" {
		t.Errorf("expected no output, got %q", got)
	}
}

func TestConditionMustBeTrueOrFalse(t *testing.T) {
	expectFails(t, "if 1 then\n  print 1\nend if", "condition has to be true or false")
	expectFails(t, "loop until \"x\"\n  print 1\nend loop", "condition has to be true or false")
}

func TestLoop(t *testing.T) {
	source := `num i = 0
loop until i is equal to 3
  print i
  i = i + 1
end loop`
	if got := run(t, source); got != "0\n1\n2\n" {
		t.Errorf("got %q", got)
	}
}

func TestLoopThatNeverRuns(t *testing.T) {
	if got := run(t, "loop until true\n  print \"no\"\nend loop"); got != "" {
		t.Errorf("expected no output, got %q", got)
	}
}

// Declaring inside a loop body has to work on every pass, not just the first.
func TestLoopBodyGetsAFreshScope(t *testing.T) {
	source := `num i = 0
loop until i is equal to 3
  num doubled = i * 2
  print doubled
  i = i + 1
end loop`
	if got := run(t, source); got != "0\n2\n4\n" {
		t.Errorf("got %q", got)
	}
}

// A runaway loop should be reported, not left to hang the machine.
func TestRunawayLoopIsCaught(t *testing.T) {
	expectFails(t, "num i = 0\nloop until i is less than 0\n  i = i + 1\nend loop", "will the condition")
}

func TestCodeblockWithGiveBack(t *testing.T) {
	source := `codeblock add with a, b
  give back a + b
end codeblock

print add with 2, 3`
	expectPrints(t, source, "5")
}

func TestCodeblockWithoutValues(t *testing.T) {
	source := `codeblock hello
  print "hi"
end codeblock

hello`
	expectPrints(t, source, "hi")
}

func TestCodeblockArity(t *testing.T) {
	source := `codeblock add with a, b
  give back a + b
end codeblock

print add with 1`
	expectFails(t, source, "needs 2 values, but was given 1")
}

func TestCodeblockCanSeeOuterVariables(t *testing.T) {
	source := `num bonus = 10
codeblock boosted with n
  give back n + bonus
end codeblock

print boosted with 5`
	expectPrints(t, source, "15")
}

// Values handed to a codeblock are copies; changing them must not leak out.
func TestArgumentsDoNotLeak(t *testing.T) {
	source := `num n = 1
codeblock bump with value
  value = value + 100
  give back value
end codeblock

print bump with n
print n`
	if got := run(t, source); got != "101\n1\n" {
		t.Errorf("got %q", got)
	}
}

func TestCodeblockCanChangeOuterVariable(t *testing.T) {
	source := `num total = 0
codeblock addToTotal with n
  total = total + n
end codeblock

addToTotal with 5
addToTotal with 7
print total`
	expectPrints(t, source, "12")
}

func TestRecursion(t *testing.T) {
	source := `codeblock factorial with n
  if n is less than 2 then
    give back 1
  end if
  give back n * factorial with n - 1
end codeblock

print factorial with 5`
	expectPrints(t, source, "120")
}

func TestGiveBackStopsTheCodeblock(t *testing.T) {
	source := `codeblock first
  give back 1
  print "never"
end codeblock

print first`
	expectPrints(t, source, "1")
}

// A "give back" inside a loop inside a codeblock has to unwind both.
func TestGiveBackFromInsideALoop(t *testing.T) {
	source := `codeblock firstOver with limit
  num i = 0
  loop until i is greater than 100
    if i is greater than limit then
      give back i
    end if
    i = i + 1
  end loop
  give back -1
end codeblock

print firstOver with 3`
	expectPrints(t, source, "4")
}

func TestEndlessRecursionIsCaught(t *testing.T) {
	source := `codeblock forever with n
  give back forever with n
end codeblock

print forever with 1`
	expectFails(t, source, "running itself forever")
}

func TestCallingSomethingThatIsNotACodeblock(t *testing.T) {
	expectFails(t, "num x = 1\nprint x with 2", "not a codeblock")
}

func TestPrintAliasesAllWork(t *testing.T) {
	if got := run(t, "print 1\nout 2\noutput 3"); got != "1\n2\n3\n" {
		t.Errorf("got %q", got)
	}
}

func TestErrorsCarryLineNumbers(t *testing.T) {
	expectFails(t, "num a = 1\nnum b = 2\nprint c", "line 3:")
}

func TestCommentsDoNotAffectOutput(t *testing.T) {
	source := `-- adds two numbers
num a = 1 -- first
num b = 2 -- second
print a + b`
	expectPrints(t, source, "3")
}

func TestSessionKeepsStateBetweenEntries(t *testing.T) {
	var out strings.Builder
	session := interpreter.NewSession(&out)

	if err := session.Run("num x = 41"); err != nil {
		t.Fatalf("first entry failed: %v", err)
	}
	if err := session.Run("x = x + 1"); err != nil {
		t.Fatalf("second entry failed: %v", err)
	}

	value, err := session.Eval("x")
	if err != nil {
		t.Fatalf("third entry failed: %v", err)
	}
	if value.Inspect() != "42" {
		t.Errorf("session lost track of x, got %q", value.Inspect())
	}
}

// One program exercising most of the language at once, as a guard against
// pieces that pass in isolation but do not fit together.
func TestEverythingTogether(t *testing.T) {
	source := `codeblock grade with score
  if score is greater than 89 then
    give back "A"
  otherwise if score is greater than 79 then
    give back "B"
  otherwise if score is greater than 69 then
    give back "C"
  otherwise
    give back "F"
  end if
end codeblock

num scores = 4
num highest = 0
num total = 0
num i = 0

loop until i is equal to scores
  num score = 60 + i * 11
  total = total + score
  if score is greater than highest then
    highest = score
  end if
  print "score " + score + " earns a " + grade with score
  i = i + 1
end loop

print "highest was " + highest
print "average was " + total / scores`

	want := strings.Join([]string{
		"score 60 earns a F",
		"score 71 earns a C",
		"score 82 earns a B",
		"score 93 earns a A",
		"highest was 93",
		"average was 76.5",
	}, "\n") + "\n"

	if got := run(t, source); got != want {
		t.Errorf("combined program printed:\n%s\nwant:\n%s", got, want)
	}
}
