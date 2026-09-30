package object

// Environment maps names to values, with an outer scope for codeblock bodies
// to fall back on. Alongside each value it remembers the type the variable was
// declared with, which is what lets Pharos catch "num score = 1" followed by
// `score = "one"`.
type Environment struct {
	store map[string]Object
	types map[string]Type
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: map[string]Object{}, types: map[string]Type{}}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	if obj, ok := e.store[name]; ok {
		return obj, true
	}
	if e.outer != nil {
		return e.outer.Get(name)
	}
	return nil, false
}

// DeclaredHere reports whether the name was declared in this scope, so that a
// repeated declaration can be flagged without shadowing complaints.
func (e *Environment) DeclaredHere(name string) bool {
	_, ok := e.store[name]
	return ok
}

// Declare introduces a name in this scope. An empty declaredType means the
// name is unconstrained, which is how codeblock parameters work.
func (e *Environment) Declare(name string, declaredType Type, val Object) {
	e.store[name] = val
	if declaredType != "" {
		e.types[name] = declaredType
	}
}

// Assign updates an existing name in whichever scope declared it, and reports
// false if the name was never declared.
func (e *Environment) Assign(name string, val Object) bool {
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return true
	}
	if e.outer != nil {
		return e.outer.Assign(name, val)
	}
	return false
}

// DeclaredType returns the type a name was declared with, if it was
// constrained to one.
func (e *Environment) DeclaredType(name string) (Type, bool) {
	if _, ok := e.store[name]; ok {
		t, constrained := e.types[name]
		return t, constrained
	}
	if e.outer != nil {
		return e.outer.DeclaredType(name)
	}
	return "", false
}
