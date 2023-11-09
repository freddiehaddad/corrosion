package lexer

import "github.com/freddiehaddad/corrosion/pkg/token"

type Lexer struct {
	input        string
	ch           byte
	position     int
	readPosition int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.readCharacter()
	switch l.ch {
	// delimiters
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	// operators
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.FSLASH, l.ch)
	case '<':
		tok = newToken(token.LESS_THAN, l.ch)
	case '>':
		tok = newToken(token.GREATER_THAN, l.ch)

	// end of input
	case token.EOF_VALUE:
		tok = newToken(token.EOF, l.ch)

	// errors
	default:
		tok = newToken(token.UNKNOWN, l.ch)
	}

	return tok
}

func (l *Lexer) readCharacter() {
	if l.eof() {
		l.ch = token.EOF_VALUE
		return
	}

	l.ch = l.input[l.readPosition]
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekCharacter() byte {
	if l.eof() {
		return token.EOF_VALUE
	}

	return l.input[l.readPosition]
}
func (l *Lexer) eof() bool {
	return l.readPosition == len(l.input)
}
