package ast

import (
	"mooshy/token"
	"testing"
)

func TestString(t *testing.T) {
	name := &Identifier{Token: token.Token{Type: token.IDENT, Literal: "x"}, Value: "x"}
	value := &Identifier{Token: token.Token{Type: token.IDENT, Literal: "y"}, Value: "y"}
	letStmt := &LetStatement{Token: token.Token{Type: token.LET, Literal: "let"}, Name: name, Value: value}

	var statements []Statement
	statements = append(statements, letStmt)
	p := &Program{Statements: statements}

	if p.String() != "let x = y;" {
		t.Errorf("program.String() wrong. got=%q", p.String())

	}

}
