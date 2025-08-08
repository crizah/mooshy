package ast

import (
	"bytes"
	"mooshy/token"
	"strconv"
)

type Node interface {
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var output bytes.Buffer
	for _, st := range p.Statements {
		output.WriteString(st.String())

	}
	return output.String()
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

func (ls *LetStatement) String() string {
	// let x = 5;
	var output bytes.Buffer

	output.WriteString(ls.TokenLiteral() + " ")
	output.WriteString(ls.Name.String() + " = ")
	if ls.Value != nil {
		output.WriteString(ls.Value.String())
	}
	output.WriteString(";")

	return output.String()

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

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) ExpressionNode() {}
func (il *IntegerLiteral) String() string {

	var output bytes.Buffer
	// convert int64 to string
	x := strconv.FormatInt(il.Value, 10)

	output.WriteString(x)
	return output.String()

}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (id *Identifier) String() string {
	var output bytes.Buffer
	// convert int64 to string

	output.WriteString(id.Value)
	return output.String()

}

func (i *Identifier) TokenLiteral() string {

	return i.Token.Literal

}
func (i *Identifier) ExpressionNode() {}

type ReturnStatement struct {
	Token token.Token
	Value *Identifier //// NEED TO BE ABLE TO RETURN AN INFIXeXPRESSION. nees to be expression
	// Value *Expression
}

func (rt *ReturnStatement) TokenLiteral() string {
	return rt.Token.Literal
}
func (rt *ReturnStatement) StatementNode() {}

func (rs *ReturnStatement) String() string {
	var output bytes.Buffer
	// return x;
	output.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		output.WriteString(rs.Value.Value)

	}
	output.WriteString(";")
	return output.String()

}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

// func (ex *ExpressionStatement.Expression) ExpressionNode() {}

type PrefixExpression struct {
	// "!5"
	Token    token.Token // !
	Operator string      // "!"
	Right    Expression  // expresseion to the right of the sign, 5
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
	Right    Expression
}

type BoolExpression struct {
	Token token.Token
	Value bool
}

func (be *BoolExpression) String() string {
	return be.Token.Literal
}

func (be *BoolExpression) TokenLiteral() string {
	return be.Token.Literal
}

func (be *BoolExpression) ExpressionNode() {}
func (in *InfixExpression) String() string {
	var output bytes.Buffer
	// convert int64 to string
	// x := strconv.FormatInt(in.Value, 10)

	output.WriteString("(" + in.Left.String() + in.Operator + in.Right.String() + ")")
	return output.String()

}

func (in *InfixExpression) TokenLiteral() string {
	return in.Token.Literal
}

func (in *InfixExpression) ExpressionNode() {}
func (pe *PrefixExpression) String() string {
	var output bytes.Buffer

	output.WriteString("(" + pe.Operator + pe.Right.String() + ")")
	return output.String()

}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) ExpressionNode() {}

func (ex *ExpressionStatement) TokenLiteral() string {
	return ex.Token.Literal
}

func (ex *ExpressionStatement) StatementNode() {}

func (ex *ExpressionStatement) String() string {
	if ex.Expression != nil {
		return ex.Expression.String()
	}
	return ""
}

type IfExpression struct {
	// if *Condition* { *Consequence* } else { *Alternative* }
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

type FunctionLiteral struct {
	Token     token.Token
	Parameter []*Identifier
	Body      *BlockStatement
}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) ExpressionNode() {}
func (fl *FunctionLiteral) String() string {
	// func(Parameter) { Body }
	var output bytes.Buffer
	output.WriteString("func ( ")
	for i, id := range fl.Parameter {
		if i == len(fl.Parameter)-1 {
			output.WriteString(id.String() + " ){")
		}
		output.WriteString(id.String() + " , ")

	}
	output.WriteString(fl.Body.String() + " }")
	return output.String()

}

func (bs *BlockStatement) String() string {
	var output bytes.Buffer

	for _, s := range bs.Statements {
		output.WriteString(s.String())
	}
	return output.String()
}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal

}
func (bs *BlockStatement) StatementNode() {}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) ExpressionNode() {}

func (ie *IfExpression) String() string {
	var output bytes.Buffer

	output.WriteString("if " + ie.Condition.String() + "{ " + ie.Consequence.String() + " }")
	if ie.Alternative != nil {
		output.WriteString("else { " + ie.Alternative.String() + " }")
	}

	return output.String()

}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) String() string {
	var output bytes.Buffer
	output.WriteString(ce.Function.String() + "(")
	for i, arg := range ce.Arguments {
		if i == len(ce.Arguments)-1 {
			output.WriteString(arg.String() + ")")
		} else {
			output.WriteString(arg.String() + ", ")

		}
	}

	return output.String()

}

func (ce *CallExpression) ExpressionNode() {}
func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}
