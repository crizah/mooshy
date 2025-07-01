package lexer

import "mooshy/token"

type Lexer struct {
	input        string
	prevPosition int  // index of char thats just been read, need to keep acount of this for parser
	currPosition int  // index of next char to read
	ch           byte // char being processed

}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() { // defined like this means we can call it like l.readChar() and it will work on l's data
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
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case 0:
		// tok = newToken(token.EOF, "") cant do this cuz newToken changes the byte to string
		tok.Literal = ""
		tok.Type = token.EOF
	}

	l.readChar() // when we call nextToken again, the pointers are updated
	return tok
}
