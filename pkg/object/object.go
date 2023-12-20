package object

import (
	"fmt"
	"strings"

	"github.com/freddiehaddad/corrosion/pkg/ast"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	FUNCTION_OBJ = "FUNCTION"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	NULL_OBJ     = "NULL"
)

type Function struct {
	Body       *ast.BlockStatement
	Env        *Environment
	Parameters []ast.Identifier
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var sb strings.Builder

	sb.WriteString("func")

	sb.WriteByte('(')
	sep := ""
	for _, p := range f.Parameters {
		sb.WriteString(sep)
		sb.WriteString(p.String())
		sep = ", "
	}
	sb.WriteByte(')')
	sb.WriteString(" {")
	sb.WriteString(f.Body.String())
	sb.WriteString(" }")

	return sb.String()
}

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType { return RETURN_OBJ }
func (r *Return) Inspect() string  { return r.Value.Inspect() }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Error struct {
	Value string
}

func (e *Error) Inspect() string  { return e.Value }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Null struct {
	Value *interface{}
}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }
