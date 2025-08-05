package parser

import (
	"mooshy/ast"
	"mooshy/lexer"
	"mooshy/token"
	"strconv"
)

const (
	_           int = iota // increments, skips 0 for prioroty of an expression
	LOWEST                 // strings, numbers
	EQUALS                 // ==
	LESSGREATER            // > or <
	SUM                    // + or -
	PRODUCT                // * or /
	PREFIX                 // -X or !X
	CALL                   // myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LESSER:   LESSGREATER,
	token.GREATER:  LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.DIVIDE:   PRODUCT,
	token.MULTIPLY: PRODUCT,
}

type Parser struct {
	lexer            *lexer.Lexer
	currToken        token.Token
	peekToken        token.Token
	prefixParseFuncs map[token.TokenType]PrefixParseFunc
	infixParseFuncs  map[token.TokenType]InfixParseFunc
}

type (
	PrefixParseFunc func() ast.Expression               // doesn't need to a function
	InfixParseFunc  func(ast.Expression) ast.Expression // the input is the expression on the left  , x+y, so its x
	// *ast.Identifier

)

// func (p *Parser) helper(y ast.Identifier) ast.Expression {
// 	return &y

// }

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currPrecedence() int {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) putPrefix(tok token.TokenType, pre PrefixParseFunc) {
	p.prefixParseFuncs[tok] = pre
	// exp := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	// // need to turn a function that returns ast.Expression type
	// // i want to return a pointer of type *ast.Identifier
	// p.prefixParseFuncs[p.currToken.Type] = exp

}

func (p *Parser) putInfix(tok token.TokenType, in InfixParseFunc) {
	p.infixParseFuncs[tok] = in
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	s := p.currToken.Literal
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return &ast.IntegerLiteral{Token: p.currToken, Value: val}
}

func (p *Parser) parseIdentifier() ast.Expression {
	// fmt.Printf("Type of Identifier: %s, Literal of identifier: %s\n", p.currToken.Type, p.currToken.Literal)

	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	// Result should be an integerLiteral() expression
	yeah := &ast.PrefixExpression{Token: p.currToken, Operator: p.currToken.Literal}
	// p.nextToken()
	p.nextToken()
	result := p.parseExpression(PREFIX)
	yeah.Right = result
	// this yeah needs to be smt := &ast.ExpressionStatement{Expreseeion: yeah}

	return yeah
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	// Left wont always be a number, can be any expression, so any of the CONSTS
	// 5 + 10
	yeah := &ast.InfixExpression{Token: p.currToken, Operator: p.currToken.Literal, Left: left}
	curr := p.currPrecedence() // of 5
	p.nextToken()
	// curr = +, peek = 10

	result := p.parseExpression(curr) // the input should be the precedance of left
	yeah.Right = result
	// peek := p.peekPrecedence()
	// right := p.parseExpression(peek)

	return yeah

}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.prefixParseFuncs = make(map[token.TokenType]PrefixParseFunc)
	p.nextToken()
	p.nextToken()

	p.putPrefix(token.IDENT, p.parseIdentifier)

	p.putPrefix(token.INT, p.parseIntegerLiteral)
	p.putPrefix(token.NOT, p.parsePrefixExpression) // can be ! or - or +
	p.putPrefix(token.MINUS, p.parsePrefixExpression)

	p.infixParseFuncs = make(map[token.TokenType]InfixParseFunc)

	p.putInfix(token.PLUS, p.parseInfixExpression)
	p.putInfix(token.MINUS, p.parseInfixExpression)
	p.putInfix(token.MULTIPLY, p.parseInfixExpression)
	p.putInfix(token.DIVIDE, p.parseInfixExpression)
	p.putInfix(token.LESSER, p.parseInfixExpression)
	p.putInfix(token.GREATER, p.parseInfixExpression)
	p.putInfix(token.NOT_EQ, p.parseInfixExpression)
	p.putInfix(token.EQ, p.parseInfixExpression)

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currToken.Type != token.EOF {
		st := p.parseStatement()
		if st != nil {
			program.Statements = append(program.Statements, st)
		}
		p.nextToken()

	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}

}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}

	stmt.Expression = p.parseExpression(LOWEST) // this needs to be PrefixExpression
	// needs to be followed by an IntegerLiteral

	if p.expectedPeek(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt

}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// curr = +
	prefix := p.prefixParseFuncs[p.currToken.Type]
	if prefix == nil {
		return nil
	}

	leftExp := prefix() // this was needed when it was a function. but we have changed that
	for !p.expectedPeek(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken} // LET

	if !p.expectedPeek(token.IDENT) { // cheking for syntax. let should always be followd by an IDENT
		// this also does nextToken()
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectedPeek(token.ASSIGN) { // should ALWAYS BE FOLLOWD BY AN =
		return nil
	}

	// TODO: We're skipping the expressions until we
	// encounter a semicolon

	// stmt.Value should be the return after the equal sign

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// currently only works for return x;
	stmt := &ast.ReturnStatement{Token: p.currToken} // can be IDENT, STRING, FUNC
	if p.peekToken.Type != token.SEMICOLON && p.peekToken.Type != token.RBRACE {
		p.nextToken()
	}

	stmt.Value = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectedPeek(token.SEMICOLON) {
		return nil
	}
	p.nextToken()
	return stmt

}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.currToken.Type == t
}

//	func (p *Parser) peekTokenIs(t token.TokenType) bool {
//		return p.peekToken.Type == t
//	}
func (p *Parser) expectedPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		return false
	}
}
