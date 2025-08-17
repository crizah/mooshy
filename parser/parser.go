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
	token.LPAREN:   CALL,
	token.LBLOCK:   CALL,
	token.ASSIGN:   CALL,
}

type Parser struct {
	lexer            *lexer.Lexer
	currToken        token.Token
	peekToken        token.Token
	prefixParseFuncs map[token.TokenType]PrefixParseFunc
	infixParseFuncs  map[token.TokenType]InfixParseFunc
	Errors           []string
}

type (
	PrefixParseFunc func() ast.Expression               // doesn't need to a function
	InfixParseFunc  func(ast.Expression) ast.Expression // the input is the expression on the left  , x+y, so its x
	// *ast.Identifier

)

// func (p *Parser) helper(y ast.Identifier) ast.Expression {
// 	return &y

// }

// NEED TO IMPLEMENT STRING AS A LITERAL yeah done

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

func (p *Parser) parseStringLiteral() ast.Expression {
	s := &ast.StringLiteral{Token: p.currToken}
	s.Value = p.currToken.Literal
	return s
}

func (p *Parser) parseIdentifier() ast.Expression {

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
	curr := p.currPrecedence()
	p.nextToken()

	result := p.parseExpression(curr) // the input should be the precedance of left
	yeah.Right = result
	// peek := p.peekPrecedence()
	// right := p.parseExpression(peek)

	return yeah

}

func (p *Parser) parseBooleanExpression() ast.Expression {
	if p.currToken.Type == token.TRUE {
		return &ast.BoolExpression{Token: p.currToken, Value: true}
	}
	return &ast.BoolExpression{Token: p.currToken, Value: false}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectedPeek(token.RPAREN) {
		e := "Expected closing parenthesis"
		p.Errors = append(p.Errors, e)
		return nil
	}
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	result := &ast.IfExpression{Token: p.currToken}
	if !p.expectedPeek(token.LPAREN) {
		e := "Expected opening parenthesis"
		p.Errors = append(p.Errors, e)
		return nil
	}

	p.nextToken()
	condition := p.parseExpression(LOWEST)
	result.Condition = condition
	if !p.expectedPeek(token.RPAREN) {
		e := "Expected closing parenthesis"
		p.Errors = append(p.Errors, e)
		return nil
	}

	if !p.expectedPeek(token.LBRACE) {
		e := "Expected opening bracket"
		p.Errors = append(p.Errors, e)
		return nil
	}

	consequence := p.parseBlockStatement()
	result.Consequence = consequence
	p.nextToken()
	if p.curTokenIs(token.ELSE) {
		if !p.expectedPeek(token.LBRACE) {
			e := "Expected opening bracket"
			p.Errors = append(p.Errors, e)
			return nil
		}

		alternative := p.parseBlockStatement()
		result.Alternative = alternative
	}

	return result

}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	bstmt := &ast.BlockStatement{Token: p.currToken}
	bstmt.Statements = []ast.Statement{}
	p.nextToken()
	for !p.curTokenIs(token.EOF) && !p.curTokenIs(token.RBRACE) {
		stmt := p.parseStatement()
		if stmt != nil {
			bstmt.Statements = append(bstmt.Statements, stmt)
		}

		p.nextToken()
	}
	return bstmt

}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	fl := &ast.FunctionLiteral{Token: p.currToken}
	if !p.expectedPeek(token.LPAREN) {
		e := "Expected opening parenthesis"
		p.Errors = append(p.Errors, e)
		return nil
	}
	var par []*ast.Identifier
	p.nextToken()
	for !p.curTokenIs(token.RPAREN) {

		if p.curTokenIs(token.COMMA) {
			p.nextToken()
		}
		curr := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		par = append(par, curr)
		p.nextToken()

	}

	fl.Parameter = par

	if !p.expectedPeek(token.LBRACE) {
		e := "Expected opening brace"
		p.Errors = append(p.Errors, e)
		return nil
	}
	// p.nextToken()

	bs := p.parseBlockStatement()
	fl.Body = bs

	return fl
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	s := &ast.CallExpression{Token: p.currToken, Function: function}
	// curr = (
	s.Arguments = p.parseArguments()
	return s
}

func (p *Parser) parseReAssignExpression(ident ast.Expression) ast.Expression {
	// x= 12
	// curr =
	s := &ast.ReAssignExpression{Token: p.currToken, Name: ident}
	p.nextToken()
	s.Value = p.parseExpression(LOWEST) // this might not work

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
		if p.curTokenIs(token.EOF) {
			e := ("expected semicolon ")
			p.Errors = append(p.Errors, e)
			return nil

		}
	}
	// p.nextToken()

	return s
}

func (p *Parser) parseIndexExpression(arr ast.Expression) ast.Expression {
	s := &ast.IndexExpression{Token: p.currToken, Name: arr}
	// curr = [
	p.nextToken()
	if !p.curTokenIs(token.INT) {
		e := "Expected Integer Value"
		p.Errors = append(p.Errors, e)
		return nil
	}

	val, err := strconv.ParseInt(p.currToken.Literal, 10, 64)
	if err != nil {
		panic(err)
	}

	s.Index = val
	p.nextToken()
	if !p.curTokenIs(token.RBLOCK) {
		e := "Expected closing ]"
		p.Errors = append(p.Errors, e)
		return nil

	}
	p.nextToken()
	return s

}

func (p *Parser) parseArrayExpression() ast.Expression {

	arr := &ast.ArrayExpression{Token: p.currToken}
	val := p.parseArguments()
	if len(val) != 0 {
		arr.Value = val
	}

	// otherwise stays nil for edge case handling

	return arr

}

func (p *Parser) parseArguments() []ast.Expression {
	args := []ast.Expression{}
	// add(x, y, x)

	// curr = (

	switch p.currToken.Literal {
	case token.LPAREN:
		p.nextToken()
		if p.curTokenIs(token.RPAREN) { // add() curr = )
			p.nextToken() // EOF
			return args
		}

		args = append(args, p.parseExpression(LOWEST)) // currToken is x
		p.nextToken()

		for !p.curTokenIs(token.RPAREN) {
			if p.curTokenIs(token.COMMA) {
				p.nextToken()
			}
			args = append(args, p.parseExpression(LOWEST))
			p.nextToken()
		}

	case token.LBLOCK:
		p.nextToken()
		if p.curTokenIs(token.RBLOCK) { // add() curr = )
			p.nextToken() // EOF
			return args
		}

		args = append(args, p.parseExpression(LOWEST)) // currToken is x
		// first
		p.nextToken()

		for !p.curTokenIs(token.RBLOCK) {
			if p.curTokenIs(token.COMMA) {
				p.nextToken()
			}

			// come on man u should also have .type be able to be used outside of switch cases

			args = append(args, p.parseExpression(LOWEST))
			p.nextToken()
		}

	}

	return args

}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l}
	p.prefixParseFuncs = make(map[token.TokenType]PrefixParseFunc)
	p.infixParseFuncs = make(map[token.TokenType]InfixParseFunc)
	p.nextToken()
	p.nextToken()

	p.putPrefix(token.IDENT, p.parseIdentifier)
	p.putPrefix(token.STRING, p.parseStringLiteral)

	p.putPrefix(token.INT, p.parseIntegerLiteral) // THIS IS NOT CALLING THE FUNCTION. ITRS JUST UK MENTIONING IT
	p.putPrefix(token.NOT, p.parsePrefixExpression)
	p.putPrefix(token.MINUS, p.parsePrefixExpression)
	p.putPrefix(token.TRUE, p.parseBooleanExpression)
	p.putPrefix(token.FALSE, p.parseBooleanExpression)

	p.putInfix(token.PLUS, p.parseInfixExpression)
	p.putInfix(token.MINUS, p.parseInfixExpression)
	p.putInfix(token.MULTIPLY, p.parseInfixExpression)
	p.putInfix(token.DIVIDE, p.parseInfixExpression)
	p.putInfix(token.LESSER, p.parseInfixExpression)
	p.putInfix(token.GREATER, p.parseInfixExpression)
	p.putInfix(token.NOT_EQ, p.parseInfixExpression)
	p.putInfix(token.EQ, p.parseInfixExpression)

	p.putPrefix(token.LPAREN, p.parseGroupedExpression) // look at this
	p.putPrefix(token.IF, p.parseIfExpression)
	p.putPrefix(token.FUNC, p.parseFunctionLiteral)

	p.putInfix(token.LPAREN, p.parseCallExpression) // in add(x,y), ( is at the infix position, after pasing the function
	// p.putPrefix(token.RETURN, p.parseReturn)
	p.putPrefix(token.LBLOCK, p.parseArrayExpression)

	p.putInfix(token.LBLOCK, p.parseIndexExpression)
	p.putInfix(token.ASSIGN, p.parseReAssignExpression)

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

	prefix := p.prefixParseFuncs[p.currToken.Type] // so u cant do type1 * type2. both need to be same time
	// prefix is -a in -a + b
	if prefix == nil {
		// e := ("PrefixParseFuncs not available for" + p.currToken.Literal)
		e := ("No attatched prefix function found")
		p.Errors = append(p.Errors, e)

		return nil
	}

	leftExp := prefix() // we are calling these here
	for !p.expectedPeek(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.Type] // parseInfixExpressin()
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp) // calling it here, with left = previous expression
	}
	return leftExp // returns InfixExpression type or can be PrefixExpression as well
	// which will the type of the expresion in stmt.ExpressionStatement.Expression
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currToken} // LET

	if !p.expectedPeek(token.IDENT) { // cheking for syntax. let should always be followd by an IDENT
		// this also does nextToken()
		e := ("expected IDENT")
		p.Errors = append(p.Errors, e)
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	if !p.expectedPeek(token.ASSIGN) { // should ALWAYS BE FOLLOWD BY AN =
		e := ("expected ASSIGN")
		p.Errors = append(p.Errors, e)
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST) // this might not work
	// p.nextToken()
	// if !p.curTokenIs(token.SEMICOLON) {
	// 	e := ("expected SEMICOLON")
	// 	p.Errors = append(p.Errors, e)
	// 	return nil
	// }

	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
		if p.curTokenIs(token.EOF) {
			e := ("expected semicolon")
			p.Errors = append(p.Errors, e)
			return nil

		}
	}
	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	stmt := &ast.ReturnStatement{Token: p.currToken}

	// if p.peekToken.Type != token.SEMICOLON && p.peekToken.Type != token.RBRACE {
	// 	p.nextToken()
	// }
	p.nextToken()
	if p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
		return stmt
	}

	stmt.Value = p.parseExpression(LOWEST)

	// stmt.Value = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	// p.nextToken()
	for !p.curTokenIs(token.SEMICOLON) {

		p.nextToken()
		if p.curTokenIs(token.EOF) {
			e := ("expected semicolon")
			p.Errors = append(p.Errors, e)
			return nil
		}
	}
	// p.nextToken()

	// p.nextToken()

	// if !p.curTokenIs(token.SEMICOLON) {
	// 	e := ("expected SEMICOLON")
	// 	p.Errors = append(p.Errors, e)
	// 	return nil
	// }
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
