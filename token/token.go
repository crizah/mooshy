package token

type TokenType string // only capital letter identifiers are exported

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// keywords
	LET  = "LET"
	FUNC = "FUNC"

	// operators
	ASSIGN = "="
	PLUS   = "+"

	// symbols
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	COMMA     = ","
	SEMICOLON = ";"

	// literals and identifiers
	INT      = "INT"
	VARIABLE = "VARIABLE"

	// errors
	ILLEGAL = "ILLEGAL" // token we have not defined
	EOF     = "EOF"     // signifies end of file

)
