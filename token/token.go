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
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	MULTIPLY = "*"
	DIVIDE   = "/"
	NOT      = "!"
	LESSER   = "<"
	GREATER  = ">"

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
	IDENT    = "IDENT" // used for identifiers that are not keywords

	// errors
	ILLEGAL = "ILLEGAL" // token we have not defined
	EOF     = "EOF"     // signifies end of file

)

// var keywords = map[string]TokenType{
// 	"let":    LET,
// 	"func":   FUNC,
// 	"if":     IF,
// 	"else":   ELSE,
// 	"return": RETURN,
// 	"true":   TRUE,
// 	"false":  FALSE,
// }
