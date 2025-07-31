package lexer

import (
	"mooshy/token"
)

type Lexer struct {
	input        string
	prevPosition int  // index of char thats just been read, need to keep acount of this for parser
	currPosition int  // index of next char to read
	ch           byte // char being processed

}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	// l.currPosition = 0 // starting poiunts
	// l.prevPosition = -1
	// l.ch = input[l.currPosition]
	l.readChar()
	return l
}

func (l *Lexer) readChar() { // defined like this means we can call it like l.readChar() and it will work on l's data
	// this is to update pointers of the lexer

	if l.currPosition >= len(l.input) {
		l.ch = 0 // end of input
	} else {
		l.ch = l.input[l.currPosition]
	}
	l.prevPosition = l.currPosition
	l.currPosition++
}

func newToken(operator token.TokenType, character byte) token.Token {
	tok := token.Token{Type: operator, Literal: string(character)}
	return tok
}

var keywords = map[string]token.TokenType{
	"let":    token.LET,
	"func":   token.FUNC,
	"if":     token.IF,
	"else":   token.ELSE,
	"return": token.RETURN,
	"true":   token.TRUE,
	"false":  token.FALSE,
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	switch l.ch {
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '=':
		if l.peekChar() == '=' {

			tok.Literal = "=="  // read the next character
			tok.Type = token.EQ // set the type to EQ
			l.readChar()

		} else {
			tok = newToken(token.ASSIGN, l.ch)

		}

	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.MULTIPLY, l.ch)
	case '/':
		tok = newToken(token.DIVIDE, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok.Literal = string(ch) + string(l.ch)
			tok.Type = token.NOT_EQ
		} else {
			tok = newToken(token.NOT, l.ch)
		}

	case '<':
		tok = newToken(token.LESSER, l.ch)
	case '>':
		tok = newToken(token.GREATER, l.ch)
	case '"':
		tok.Literal = l.readString() // read the string literal
		tok.Type = token.STRING
		return tok

	case 0:
		// tok = newToken(token.EOF, "") cant do this cuz newToken changes the byte to string
		//implement that zero doesnt always have to mean end of file it can be read as 019 which is just int
		tok.Literal = ""
		tok.Type = token.EOF

	default:
		if isLetter(l.ch) { // can be keyword or identifier

			tok.Literal = l.readIdentifier() // function reads the string until non-letter nd then returns its literal
			tok.Type = l.getType(tok.Literal)
			// l.readChar()
			return tok

		}
		if isWhitespace(l.ch) {
			for isWhitespace(l.ch) {
				l.readChar()
			}
			return l.NextToken()
		}
		if isNumber(l.ch) {

			tok.Literal = l.readNumber() // function reads the string until non-number and then returns its literal
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch) // if its not a letter, then its an illegal token
		}

	}

	l.readChar() // when we call nextToken again, the pointers are updated
	return tok
}

func (l *Lexer) readIdentifier() string {
	// read until its not a letter and return its literal
	s := ""

	for isLetter(l.ch) {
		s += string(l.ch) // append the character to the string
		l.readChar()      // move to the next character

	}

	return s

}

func (l *Lexer) readString() string {
	// red until next "
	s := ""
	// s += string(l.ch)
	l.readChar()
	for l.ch != '"' {
		s += string(l.ch)
		l.readChar()
	}
	// s += string(l.ch) //the closing "
	l.readChar() // read the closing "
	return s

}
func (l *Lexer) peekChar() byte {
	if l.currPosition >= len(l.input) {
		return 0 // end of input
	} else {
		return l.input[l.currPosition] // return the next character without moving the pointer
	}
}
func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'

}

func (l *Lexer) readNumber() string {
	// read until no more number and return it as string literal
	s := ""
	for isNumber(l.ch) {
		s += string(l.ch)
		l.readChar()

	}

	return s
}

func (l *Lexer) getType(literal string) token.TokenType {

	if tok, ok := keywords[literal]; ok { // map lookup with boolean check
		return tok // if its a keyword, return the keyword as a token type
	}
	return token.IDENT // otherwise, its an identifier or variable

}
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || (ch == '_')
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
