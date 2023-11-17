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
	l := &Lexer{input: input}
	l.readCharacter()
	return l
}

func newTokenByte(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func newTokenString(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

// Returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.consumeWhitespace()

	switch l.ch {
	// delimiters
	case ';':
		tok = newTokenByte(token.SEMICOLON, l.ch)

	// operators
	case '=':
		tok = newTokenByte(token.ASSIGN, l.ch)

	// end of input
	case token.EOF_VALUE:
		tok = newTokenByte(token.EOF, l.ch)

	default:
		if isAlpha(l.ch) {
			s := l.readWord()
			tt := token.LookupType(s)
			tok = newTokenString(tt, s)
			return tok
		} else if isDigit(l.ch) {
			s := l.readNumber()
			tok = newTokenString(token.INTEGER, s)
			return tok
		} else {
			tok = newTokenByte(token.ILLEGAL, l.ch)
		}
	}

	l.readCharacter()

	return tok
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

func isWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\r', '\n', '\t':
		return true
	default:
		return false
	}
}

func (l *Lexer) consumeWhitespace() {
	for isWhitespace(l.ch) {
		l.readCharacter()
	}
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

func (l *Lexer) readWord() string {
	sb := strings.Builder{}

	for isAlpha(l.ch) {
		sb.WriteByte(l.ch)
		l.readCharacter()
	}

	return sb.String()
}

func (l *Lexer) readNumber() string {
	sb := strings.Builder{}

	for isDigit(l.ch) {
		sb.WriteByte(l.ch)
		l.readCharacter()
	}

	return sb.String()
}
