// Package evaluator walks a parsed Pharos program and runs it.
package evaluator

import (
	"fmt"
	"io"

	"github.com/abhinavtripathy/pharos/interpreter/ast"
	"github.com/abhinavtripathy/pharos/interpreter/object"
)

// Runaway programs are a fact of life when you are learning, so both loops and
// codeblock calls are capped and reported as ordinary errors rather than being
// left to hang or to crash the host.
const (
	MaxLoopPasses = 1_000_000
	MaxCallDepth  = 10_000
)

var (
	trueObj    = &object.Boolean{Value: true}
	falseObj   = &object.Boolean{Value: false}
	nothingObj = &object.Nothing{}
)

type Evaluator struct {
	out   io.Writer
	depth int
}

func New(out io.Writer) *Evaluator { return &Evaluator{out: out} }

func (e *Evaluator) Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return e.evalProgram(node, env)
	case *ast.Block:
		return e.evalBlock(node, env)

	case *ast.Declare:
		return e.evalDeclare(node, env)
	case *ast.Assign:
		return e.evalAssign(node, env)
	case *ast.Print:
		return e.evalPrint(node, env)
	case *ast.If:
		return e.evalIf(node, env)
	case *ast.Loop:
		return e.evalLoop(node, env)
	case *ast.Codeblock:
		return e.evalCodeblock(node, env)
	case *ast.GiveBack:
		return e.evalGiveBack(node, env)
	case *ast.ExpressionStatement:
		return e.Eval(node.Expression, env)

	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.TextLiteral:
		return &object.Text{Value: node.Value}
	case *ast.BooleanLiteral:
		return boolean(node.Value)
	case *ast.Identifier:
		obj, ok := env.Get(node.Value)
		if !ok {
			return errorf(node, "I have not heard of %q yet", node.Value)
		}
		// Pharos has no way to hold on to a codeblock as a value, so a
		// codeblock's name appearing in an expression can only mean "run it".
		if fn, isBlock := obj.(*object.Codeblock); isBlock {
			return e.callCodeblock(node, fn, nil)
		}
		return obj

	case *ast.Prefix:
		return e.evalPrefix(node, env)
	case *ast.Infix:
		return e.evalInfix(node, env)
	case *ast.Call:
		return e.evalCall(node, env)
	}

	return errorf(node, "I do not know how to run this")
}

func (e *Evaluator) evalProgram(node *ast.Program, env *object.Environment) object.Object {
	var result object.Object = nothingObj
	for _, stmt := range node.Statements {
		result = e.Eval(stmt, env)
		switch r := result.(type) {
		case *object.Error:
			return r
		case *object.ReturnValue:
			return r.Value
		}
	}
	return result
}

func (e *Evaluator) evalBlock(node *ast.Block, env *object.Environment) object.Object {
	var result object.Object = nothingObj
	for _, stmt := range node.Statements {
		result = e.Eval(stmt, env)
		switch result.(type) {
		case *object.Error, *object.ReturnValue:
			return result
		}
	}
	return result
}

func (e *Evaluator) evalDeclare(node *ast.Declare, env *object.Environment) object.Object {
	name := node.Name.Value
	if env.DeclaredHere(name) {
		return errorf(node, "%q has already been declared here; use %q = ... to change it", name, name)
	}

	val := e.Eval(node.Value, env)
	if object.IsError(val) {
		return val
	}

	declared := object.Type(node.DeclaredType)
	if val.Type() != declared {
		return errorf(node, "%q was declared as %s, so it cannot hold %s",
			name, declared, object.Article(val.Type()))
	}

	env.Declare(name, declared, val)
	return nothingObj
}

func (e *Evaluator) evalAssign(node *ast.Assign, env *object.Environment) object.Object {
	name := node.Name.Value
	if _, ok := env.Get(name); !ok {
		return errorf(node, "%q has not been declared yet; start with something like %q", name, "num "+name+" = 0")
	}

	val := e.Eval(node.Value, env)
	if object.IsError(val) {
		return val
	}

	if declared, constrained := env.DeclaredType(name); constrained && val.Type() != declared {
		return errorf(node, "%q was declared as %s, so it cannot hold %s",
			name, declared, object.Article(val.Type()))
	}

	env.Assign(name, val)
	return nothingObj
}

func (e *Evaluator) evalPrint(node *ast.Print, env *object.Environment) object.Object {
	val := e.Eval(node.Value, env)
	if object.IsError(val) {
		return val
	}
	fmt.Fprintln(e.out, val.Inspect())
	return nothingObj
}

func (e *Evaluator) evalIf(node *ast.If, env *object.Environment) object.Object {
	for _, branch := range node.Branches {
		met, err := e.condition(branch.Condition, env)
		if err != nil {
			return err
		}
		if met {
			return e.Eval(branch.Body, object.NewEnclosedEnvironment(env))
		}
	}
	if node.Else != nil {
		return e.Eval(node.Else, object.NewEnclosedEnvironment(env))
	}
	return nothingObj
}

// evalLoop runs the body until the condition becomes true. Each pass gets a
// fresh scope so that declaring a variable inside the body works every time
// around.
func (e *Evaluator) evalLoop(node *ast.Loop, env *object.Environment) object.Object {
	for pass := 0; ; pass++ {
		if pass >= MaxLoopPasses {
			return errorf(node, "this loop has run %d times; will the condition after %q ever be true?",
				MaxLoopPasses, "until")
		}

		met, err := e.condition(node.Condition, env)
		if err != nil {
			return err
		}
		if met {
			return nothingObj
		}

		result := e.Eval(node.Body, object.NewEnclosedEnvironment(env))
		switch result.(type) {
		case *object.Error, *object.ReturnValue:
			return result
		}
	}
}

func (e *Evaluator) evalCodeblock(node *ast.Codeblock, env *object.Environment) object.Object {
	name := node.Name.Value
	if env.DeclaredHere(name) {
		return errorf(node, "%q has already been declared here", name)
	}

	env.Declare(name, object.CodeblockType, &object.Codeblock{
		Name:       name,
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	})
	return nothingObj
}

func (e *Evaluator) evalGiveBack(node *ast.GiveBack, env *object.Environment) object.Object {
	if node.Value == nil {
		return &object.ReturnValue{Value: nothingObj}
	}
	val := e.Eval(node.Value, env)
	if object.IsError(val) {
		return val
	}
	return &object.ReturnValue{Value: val}
}

func (e *Evaluator) evalPrefix(node *ast.Prefix, env *object.Environment) object.Object {
	right := e.Eval(node.Right, env)
	if object.IsError(right) {
		return right
	}

	switch node.Operator {
	case "not":
		b, ok := right.(*object.Boolean)
		if !ok {
			return errorf(node, "%q needs true or false after it, but found %s", "not", object.Article(right.Type()))
		}
		return boolean(!b.Value)

	case "-":
		n, ok := right.(*object.Number)
		if !ok {
			return errorf(node, "only a num can be negative, but this is %s", object.Article(right.Type()))
		}
		return &object.Number{Value: -n.Value}
	}

	return errorf(node, "I do not know the operator %q", node.Operator)
}

func (e *Evaluator) evalInfix(node *ast.Infix, env *object.Environment) object.Object {
	// "and" and "or" stop early, so the right side is only run when it matters.
	if node.Operator == "and" || node.Operator == "or" {
		return e.evalLogical(node, env)
	}

	left := e.Eval(node.Left, env)
	if object.IsError(left) {
		return left
	}
	right := e.Eval(node.Right, env)
	if object.IsError(right) {
		return right
	}

	switch node.Operator {
	case "+":
		return e.evalPlus(node, left, right)
	case "-", "*", "/":
		return e.evalArithmetic(node, left, right)
	case "is greater than", "is less than":
		return e.evalOrdering(node, left, right)
	case "is equal to", "is not equal to":
		return e.evalEquality(node, left, right)
	}

	return errorf(node, "I do not know the operator %q", node.Operator)
}

func (e *Evaluator) evalLogical(node *ast.Infix, env *object.Environment) object.Object {
	left := e.Eval(node.Left, env)
	if object.IsError(left) {
		return left
	}
	lb, ok := left.(*object.Boolean)
	if !ok {
		return errorf(node, "%q needs true or false on its left, but found %s",
			node.Operator, object.Article(left.Type()))
	}

	if node.Operator == "and" && !lb.Value {
		return falseObj
	}
	if node.Operator == "or" && lb.Value {
		return trueObj
	}

	right := e.Eval(node.Right, env)
	if object.IsError(right) {
		return right
	}
	rb, ok := right.(*object.Boolean)
	if !ok {
		return errorf(node, "%q needs true or false on its right, but found %s",
			node.Operator, object.Article(right.Type()))
	}
	return boolean(rb.Value)
}

// evalPlus adds numbers, and joins text to whatever is beside it so that
// `print "score: " + score` behaves the way it reads.
func (e *Evaluator) evalPlus(node *ast.Infix, left, right object.Object) object.Object {
	if l, ok := left.(*object.Number); ok {
		if r, ok := right.(*object.Number); ok {
			return &object.Number{Value: l.Value + r.Value}
		}
	}

	_, leftIsText := left.(*object.Text)
	_, rightIsText := right.(*object.Text)
	if (leftIsText || rightIsText) && joinable(left) && joinable(right) {
		return &object.Text{Value: left.Inspect() + right.Inspect()}
	}

	return errorf(node, "cannot add %s to %s", object.Article(right.Type()), object.Article(left.Type()))
}

func joinable(obj object.Object) bool {
	switch obj.(type) {
	case *object.Text, *object.Number, *object.Boolean:
		return true
	}
	return false
}

func (e *Evaluator) evalArithmetic(node *ast.Infix, left, right object.Object) object.Object {
	l, lok := left.(*object.Number)
	r, rok := right.(*object.Number)
	if !lok || !rok {
		return errorf(node, "%q only works with nums, but found %s and %s",
			node.Operator, object.Article(left.Type()), object.Article(right.Type()))
	}

	switch node.Operator {
	case "-":
		return &object.Number{Value: l.Value - r.Value}
	case "*":
		return &object.Number{Value: l.Value * r.Value}
	case "/":
		if r.Value == 0 {
			return errorf(node, "cannot divide by zero")
		}
		return &object.Number{Value: l.Value / r.Value}
	}

	return errorf(node, "I do not know the operator %q", node.Operator)
}

func (e *Evaluator) evalOrdering(node *ast.Infix, left, right object.Object) object.Object {
	l, lok := left.(*object.Number)
	r, rok := right.(*object.Number)
	if !lok || !rok {
		return errorf(node, "%q only works with nums, but found %s and %s",
			node.Operator, object.Article(left.Type()), object.Article(right.Type()))
	}

	if node.Operator == "is greater than" {
		return boolean(l.Value > r.Value)
	}
	return boolean(l.Value < r.Value)
}

func (e *Evaluator) evalEquality(node *ast.Infix, left, right object.Object) object.Object {
	if left.Type() != right.Type() {
		return errorf(node, "cannot compare %s with %s",
			object.Article(left.Type()), object.Article(right.Type()))
	}

	var same bool
	switch l := left.(type) {
	case *object.Number:
		same = l.Value == right.(*object.Number).Value
	case *object.Text:
		same = l.Value == right.(*object.Text).Value
	case *object.Boolean:
		same = l.Value == right.(*object.Boolean).Value
	default:
		return errorf(node, "%s cannot be compared", object.Article(left.Type()))
	}

	if node.Operator == "is not equal to" {
		return boolean(!same)
	}
	return boolean(same)
}

func (e *Evaluator) evalCall(node *ast.Call, env *object.Environment) object.Object {
	obj, ok := env.Get(node.Function.Value)
	if !ok {
		return errorf(node, "I have not heard of %q yet", node.Function.Value)
	}

	fn, ok := obj.(*object.Codeblock)
	if !ok {
		// A bare name is parsed as a call with no values so that codeblocks
		// taking nothing can be run by name. When the name holds an ordinary
		// value instead, it simply is that value, which is what makes typing
		// `score` on its own useful in the prompt.
		if len(node.Arguments) == 0 {
			return obj
		}
		return errorf(node, "%q is %s, not a codeblock, so it cannot be given values",
			node.Function.Value, object.Article(obj.Type()))
	}

	args := make([]object.Object, 0, len(node.Arguments))
	for _, argNode := range node.Arguments {
		arg := e.Eval(argNode, env)
		if object.IsError(arg) {
			return arg
		}
		args = append(args, arg)
	}

	return e.callCodeblock(node, fn, args)
}

// callCodeblock runs fn with args already evaluated. Parameters are declared
// without a type, so they take whatever they are handed.
func (e *Evaluator) callCodeblock(node ast.Node, fn *object.Codeblock, args []object.Object) object.Object {
	if len(args) != len(fn.Parameters) {
		return errorf(node, "%q needs %s, but was given %d",
			fn.Name, plural(len(fn.Parameters), "value"), len(args))
	}

	if e.depth >= MaxCallDepth {
		return errorf(node, "%q has been called %d times without finishing; is it running itself forever?",
			fn.Name, MaxCallDepth)
	}
	e.depth++
	defer func() { e.depth-- }()

	inner := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		inner.Declare(param.Value, "", args[i])
	}

	result := e.Eval(fn.Body, inner)
	if object.IsError(result) {
		return result
	}
	if given, ok := result.(*object.ReturnValue); ok {
		return given.Value
	}
	return nothingObj
}

// condition evaluates something used as a yes/no test, insisting it really is
// true or false rather than quietly treating other values as truthy.
func (e *Evaluator) condition(node ast.Expression, env *object.Environment) (bool, object.Object) {
	val := e.Eval(node, env)
	if object.IsError(val) {
		return false, val
	}
	b, ok := val.(*object.Boolean)
	if !ok {
		return false, errorf(node, "a condition has to be true or false, but this one is %s",
			object.Article(val.Type()))
	}
	return b.Value, nil
}

func boolean(b bool) *object.Boolean {
	if b {
		return trueObj
	}
	return falseObj
}

func errorf(node ast.Node, format string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...), Line: node.Token().Line}
}

func plural(n int, word string) string {
	if n == 1 {
		return fmt.Sprintf("%d %s", n, word)
	}
	return fmt.Sprintf("%d %ss", n, word)
}
