package object

import (
	"fmt"
	// "mooshy/token"
)

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOL_OBJ    = "BOOL"
	STRING_OBJ  = "STRING"
	NULL_OBJ    = "NULL"
	RETURN_OBJ  = "RETURN"
	ERROR_OBJ   = "ERROR"
)

type Error struct {
	Msg string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("error: %s", e.Msg)
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Bool struct {
	Value bool
}
type String struct {
	// this is not string. it should retunr Token.String. thats why not recoignied as StringObject yet
	// "string"
	Value string
}

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType {
	return RETURN_OBJ
}

func (r *Return) Inspect() string {
	return fmt.Sprintf("RETURN %s", r.Value.Inspect())
}

type Null struct{}

func (nl *Null) Inspect() string {
	return "null"
}
func (nl *Null) Type() ObjectType {
	return NULL_OBJ
}
func (str *String) Inspect() string {
	return fmt.Sprint(str.Value)
}

func (str *String) Type() ObjectType {
	return STRING_OBJ
}

func (b *Bool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Bool) Type() ObjectType {
	return BOOL_OBJ
}
func (in *Integer) Inspect() string {
	return fmt.Sprintf("%d", in.Value)
}

func (in *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
