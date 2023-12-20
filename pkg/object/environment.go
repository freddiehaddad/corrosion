// The Environment serves as a data store for variable and function
// declarations.
package object

import (
	"fmt"
)

// Environment represents the state of the environment, both globally and
// scoped environments (i.e. within scoped blocks and function calls).
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates the top level (global) environment.  Additional
// environments for scoping (i.e. block statements, functions) should use
// NewScopedEnvironment.
func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store}
}

// Creates a scoped environment that is part of function calls and block
// statements.
func NewScopedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Checks the current environment (including outer scopes) for the identifier
// name and returns its value if found along with the value true.  Otherwise,
// obj is undefined and ok will be false. Always check the result of ok before
func (e *Environment) Get(name string) (obj Object, ok bool) {
	obj, ok = e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return
}

// Sets the mapping of name to value in the current environment. Returns value.
// Should be called for new declarations.  Update should be used instead for
// updating existing variables.
//
// Example:
//
//	var foo = 100;  Handled by Set
//	foo = 200;      Should be handled with Update
func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

// Update assigns the value to the existing identifier name and returns value
// and true.  If name does not exist (meaning it hasn't already been declared),
// then value is returned and false.
func (e *Environment) Update(name string, value Object) (Object, bool) {
	if _, ok := e.store[name]; ok {
		e.store[name] = value
		return value, ok
	}

	if e.outer != nil {
		return e.outer.Update(name, value)
	}

	m := fmt.Sprintf("ERROR: undefined variable %q", name)
	err := &Error{Value: m}
	return err, false
}
