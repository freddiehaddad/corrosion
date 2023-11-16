package lexer

import (
	"strings"

	"github.com/freddiehaddad/corrosion/pkg/token"
)

type Lexer struct {
	input        string
	ch           byte
	position     int
	readPosition int
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func newTokenByte(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newTokenString(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	var sb strings.Builder

	l.readCharacter()
	switch l.ch {
	// delimiters
	case ';':
		tok = newTokenByte(token.SEMICOLON, l.ch)
	case ',':
		tok = newTokenByte(token.COMMA, l.ch)
	case '(':
		tok = newTokenByte(token.LPAREN, l.ch)
	case ')':
		tok = newTokenByte(token.RPAREN, l.ch)
	case '[':
		tok = newTokenByte(token.LBRACKET, l.ch)
	case ']':
		tok = newTokenByte(token.RBRACKET, l.ch)
	case '{':
		tok = newTokenByte(token.LBRACE, l.ch)
	case '}':
		tok = newTokenByte(token.RBRACE, l.ch)

	// operators
	case '+':
		tok = newTokenByte(token.PLUS, l.ch)
	case '-':
		tok = newTokenByte(token.MINUS, l.ch)
	case '*':
		tok = newTokenByte(token.ASTERISK, l.ch)
	case '/':
		tok = newTokenByte(token.FSLASH, l.ch)
	case '<':
		tok = newTokenByte(token.LESS_THAN, l.ch)
	case '>':
		tok = newTokenByte(token.GREATER_THAN, l.ch)
	case '^':
		tok = newTokenByte(token.BITWISE_XOR, l.ch)
	case '~':
		tok = newTokenByte(token.BITWISE_NOT, l.ch)

	// extended operators
	case '=':
		sb.WriteByte(l.ch)
		if l.peekCharacter() == '=' {
			l.readCharacter()
			sb.WriteByte(l.ch)
			tok = newTokenString(token.EQUAL, sb.String())
		} else {
			tok = newTokenByte(token.ASSIGN, l.ch)
		}
	case '!':
		sb.WriteByte(l.ch)
		if l.peekCharacter() == '=' {
			l.readCharacter()
			sb.WriteByte(l.ch)
			tok = newTokenString(token.NOT_EQUAL, sb.String())
		} else {
			tok = newTokenByte(token.BANG, l.ch)
		}
	case '&':
		sb.WriteByte(l.ch)
		if l.peekCharacter() == '&' {
			l.readCharacter()
			sb.WriteByte(l.ch)
			tok = newTokenString(token.AND, sb.String())
		} else {
			tok = newTokenByte(token.BITWISE_AND, l.ch)
		}
	case '|':
		sb.WriteByte(l.ch)
		if l.peekCharacter() == '|' {
			l.readCharacter()
			sb.WriteByte(l.ch)
			tok = newTokenString(token.OR, sb.String())
		} else {
			tok = newTokenByte(token.BITWISE_OR, l.ch)
		}

	// end of input
	case token.EOF_VALUE:
		tok = newTokenByte(token.EOF, l.ch)

	// errors
	default:
		tok = newTokenByte(token.UNKNOWN, l.ch)
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
