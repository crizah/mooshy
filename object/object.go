package object

import (
	"bytes"
	"fmt"
	"mooshy/ast"
	// "mooshy/token"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOL_OBJ     = "BOOL"
	STRING_OBJ   = "STRING"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	FUNC_OBJ     = "FUNCTION"
	BUILT_IN_OBJ = "BUILT_IN"
	ARRAY_OBJ    = "ARRAY"
	FORLOOP_OBJ  = "FORLOOP"
)

type BuiltInFunction func(args ...Object) Object

type BuiltIn struct {
	Fn BuiltInFunction
}

func (b *BuiltIn) Type() ObjectType {
	return BUILT_IN_OBJ
}

func (b *BuiltIn) Inspect() string {
	return "BuiltIn Function"
}

type Error struct {
	Msg string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR %s", e.Msg)
	// return ERROR_OBJ
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Array struct {
	Value []Object
}

func (arr *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (arr *Array) Inspect() string {

	var output bytes.Buffer

	output.WriteString("[")
	for i, p := range arr.Value {
		if i == len(arr.Value)-1 {
			output.WriteString(p.Inspect() + "]")

		} else {
			output.WriteString(p.Inspect() + ", ")
		}
	}

	return output.String()

}

type Bool struct {
	Value bool
}
type String struct {
	// this is not string. it should retunr Token.String. thats why not recoignied as StringObject yet
	// "string"
	Value string
}

type ForLoop struct {
	// Name      Integer // Integer object
	// Condition Bool
	// Iterator  Integer
	// Todo      []Object

	Name      Object // Integer object
	Condition Object
	Iterator  Object
	Todo      []Object
}

func (f *ForLoop) Type() ObjectType {
	return FORLOOP_OBJ

}

func (f *ForLoop) Inspect() string {
	var output bytes.Buffer

	// func(x){ }

	output.WriteString("for(")
	output.WriteString(f.Name.Inspect() + "; ")
	output.WriteString(f.Condition.Inspect() + "; ")
	output.WriteString(f.Iterator.Inspect() + "){")
	for _, s := range f.Todo {
		output.WriteString(s.Inspect())
	}

	output.WriteString("}")
	return output.String()

	// return FUNC_OBJ

}

type Function struct {
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Enviorment
}

func (f *Function) Inspect() string {
	var output bytes.Buffer

	// func(x){ }

	output.WriteString("func(")
	for i, p := range f.Params {
		if i == len(f.Params)-1 {
			output.WriteString(p.String() + "){")

		} else {
			output.WriteString(p.String() + ", ")
		}
	}

	for _, s := range f.Body.Statements {
		output.WriteString(s.String())
	}

	output.WriteString("}")
	return output.String()

	// return FUNC_OBJ

}

func (f *Function) Type() ObjectType {
	return FUNC_OBJ

}

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType {
	return RETURN_OBJ
}

func (r *Return) Inspect() string {
	return fmt.Sprintf("RETURN %s", r.Value.Inspect())
	// return RETURN_OBJ
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
	// return STRING_OBJ

}

func (str *String) Type() ObjectType {
	return STRING_OBJ
}

func (b *Bool) Inspect() string {
	return fmt.Sprintf("%t", b.Value)

	// return BOOL_OBJ
}

func (b *Bool) Type() ObjectType {
	return BOOL_OBJ
}
func (in *Integer) Inspect() string {
	return fmt.Sprintf("%d", in.Value)
	// return INTEGER_OBJ
}

func (in *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
