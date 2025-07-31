// lexer: It will take source code as input and output the tokens that represent the source code.
package lexer

import (
	"testing"

	"mooshy/token"
)

func TestToken(t *testing.T) {
	// ` for multiline strings
	// implement with number1 as a variable name
	input := `let numberOne = 5;
let numberTwo = 10;
let add = func(x, y) {
x + y;
};
let result = add(numberOne, numberTwo);

5*/ 10< > 

`

	tests := []struct {
		expectedToken   token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "numberOne"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "numberTwo"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNC, "func"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "numberOne"},
		{token.COMMA, ","},
		{token.IDENT, "numberTwo"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input) // New needs to have a return type of lexer
	// lexer needs to have function NextToken()

	for i, x := range tests {
		tok := l.NextToken()

		if tok.Type != x.expectedToken {
			t.Fatalf("tests[%d]- tokentype wrong. expected=%q, got=%q", i, x.expectedToken, tok.Type)
		}

		if tok.Literal != x.expectedLiteral {
			t.Fatalf("tests[%d]- literal wrong. expected=%q, got=%q", i, x.expectedLiteral, tok.Literal)
		}

	}

}
