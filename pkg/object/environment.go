package object

import (
	"fmt"
)

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewScopedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	store := make(map[string]Object)
	return &Environment{store: store}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, value Object) Object {
	e.store[name] = value
	return value
}

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
