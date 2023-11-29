package object

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJ = "INTEGER"
	NULL_OBJ    = "NULL"
)

type Null struct {
	Value string
}

func (n *Null) Inspect() string  { return n.Value }
func (n *Null) Type() ObjectType { return NULL_OBJ }

type Integer struct {
	Value string
}

func (i *Integer) Inspect() string  { return i.Value }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
