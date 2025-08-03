package token

type TokenType string // only capital letter identifiers are exported to other files

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// keywords
	LET    = "LET"
	FUNC   = "FUNC"
	RETURN = "RETURN"
	IF     = "IF"
	ELSE   = "ELSE"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	// operators
	ASSIGN       = "="
	PLUS         = "+"
	MINUS        = "-"
	MULTIPLY     = "*"
	DIVIDE       = "/"
	NOT          = "!"
	LESSER       = "<"
	GREATER      = ">"
	EQ           = "=="
	NOT_EQ       = "!="
	QUOTE        = '"'
	LESSER_EQAL  = "<="
	GREATER_EQAL = ">="
	INCREMENT    = "++"
	DECREMENT    = "--"

	// symbols
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	COMMA     = ","
	SEMICOLON = ";"

	// literals and identifiers
	INT    = "INT"
	IDENT  = "IDENT"
	STRING = "STRING"

	// errors
	ILLEGAL = "ILLEGAL" // token we have not defined
	EOF     = "EOF"     // signifies end of file

)
