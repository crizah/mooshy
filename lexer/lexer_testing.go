// lexer: It will take source code as input and output the tokens that represent the source code.
package lexer

import (
	"testing"

	"mooshy/token"
)

func testToken(t *testing.T) {
	input := "(,=};"

	tests := []struct {
		expectedToken   token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.COMMA, ","},
		{token.ASSIGN, "="},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
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
