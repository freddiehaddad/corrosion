package object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	ERROR_OBJ   = "ERROR"
	NULL_OBJ    = "NULL"
)

type Boolean struct {
	Value string
}

func (b *Boolean) Inspect() string  { return b.Value }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

type Integer struct {
	Value string
}

func (i *Integer) Inspect() string  { return i.Value }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Error struct {
	Value string
}

func (e *Error) Inspect() string  { return e.Value }
func (e *Error) Type() ObjectType { return ERROR_OBJ }

type Null struct {
	Value string
}

func (n *Null) Inspect() string  { return n.Value }
func (n *Null) Type() ObjectType { return NULL_OBJ }
