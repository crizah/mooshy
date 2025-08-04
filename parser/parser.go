package parser

import (
	"mooshy/ast"
	"mooshy/lexer"
	"mooshy/token"
	"strconv"
)

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

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.prefixParseFuncs = make(map[token.TokenType]PrefixParseFunc)
	p.nextToken()
	p.nextToken()

	// p.putPrefix(token.IDENT, p.parseIdentifier)
	p.putPrefix(token.INT, p.parseIntegerLiteral)

	// "5; " p.currToken = 5, p.peektoken = ;

	// switch p.currToken.Type{
	// 	case token.IDENT:
	// 		p.putPrefix(token.IDENT, p.parseIdentifier),
	//     case token.INT:
	// 		p.putPrefix(token.INT, p.parseIdentifier),

	// }

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
	stmt.Expression = p.parseExpression(LOWEST)

	if p.expectedPeek(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt

}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFuncs[p.currToken.Type]
	if prefix == nil {
		return nil
	}

	// leftExp := prefix()  // this was needed when it was a function. but we have changed that
	return prefix()
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

const (
	_ int = iota // increments, skips 0 for prioroty of an expression
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)
