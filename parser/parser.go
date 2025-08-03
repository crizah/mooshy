package parser

import (
	"mooshy/ast"
	"mooshy/lexer"
	"mooshy/token"
)

type Parser struct {
	lexer     *lexer.Lexer
	currToken token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.nextToken()
	p.nextToken()
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
		return nil
	}

}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken} // LET

	if !p.expectedPeek(token.IDENT) { // cheking for syntax. let should always be followd by an IDENT
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
	stmt := &ast.ReturnStatement{Token: p.currToken}
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
