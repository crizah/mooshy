package ast

import "mooshy/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	// this can also use TokenLiteral()
	Node // this means that types can use TokenLiteral() (all defined in Node) as well
	StatementNode()
}

type Expression interface {
	// this too
	Node
	ExpressionNode()
}

type Program struct {
	// array of statements, each has its own TokenLiteral() and StatementNode()
	// so uk, letStatements, ifStatements, each that will have their own defined TokenLiteralsand such
	// Program itself has its own TokenLiteral() method defined below and is a node

	Statements []Statement
}

func (p *Program) TokenLiteral() string {

	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

type LetStatement struct {
	// this is a node of type Statement

	// let x = 5;
	// Token here is LET
	// we also need its literal
	// the identifier is x and we need its Token and literal as well, so Name = x.token, x.literal
	// expression is 5
	Token token.Token // LET
	Name  *Identifier // x
	Value Expression  // this is a node of type Experession, 5
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}
func (ls *LetStatement) StatementNode() {}

// type Expression struct {
// 	Token token.Token
// 	Value string
// }

// func (ex *Expression) TokenLiteral() string {
// 	return ex.Token.Literal
// }

type Identifier struct {
	// its suppposed to be node of type Expression
	// this doesnt produce a value in let x = 5;
	// but can produce value in let x = func(x,y) { x + y; };  thats why expression and not statement

	Token token.Token
	Value string //idk what the pint of this is, it can just be TokenLiteral()
}

func (i *Identifier) TokenLiteral() string {
	// here, we need to define Token.Literal() as its not included in Identifier struct
	return i.Token.Literal

}
func (i *Identifier) ExpressionNode() {}

type ReturnStatement struct {
	Token token.Token
	Value *Identifier //can be IDENT, STRING, INT, FUNC etc
}

func (rt *ReturnStatement) TokenLiteral() string {
	return rt.Token.Literal
}
func (rt *ReturnStatement) StatementNode() {}
