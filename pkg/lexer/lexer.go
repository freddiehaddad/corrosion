package lexer

import (
	"strings"

	"github.com/freddiehaddad/corrosion/pkg/token"
)

const capacity = 10

type Lexer struct {
	tokens       chan token.Token
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.tokens = make(chan token.Token, capacity)
	l.readCharacter()
	go l.generateTokens()
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
	tok := <-l.tokens
	return tok
}

func (l *Lexer) generateTokens() {
	l.consumeWhitespace()

	for l.ch != token.EOF_VALUE {
		var tok token.Token

		switch l.ch {
		// delimiters
		case ';':
			tok = newTokenByte(token.SEMICOLON, l.ch)

		// operators
		case '=':
			tok = newTokenByte(token.ASSIGN, l.ch)
		case '-':
			tok = newTokenByte(token.MINUS, l.ch)
		case '+':
			tok = newTokenByte(token.PLUS, l.ch)
		case '*':
			tok = newTokenByte(token.MULTIPLY, l.ch)

		default:
			if isAlpha(l.ch) {
				s := l.readWord()
				tt := token.LookupType(s)
				tok = newTokenString(tt, s)
			} else if isDigit(l.ch) {
				s := l.readNumber()
				tok = newTokenString(token.INTEGER, s)
			} else {
				tok = newTokenByte(token.ILLEGAL, l.ch)
			}
		}

		l.tokens <- tok

		l.readCharacter()
		l.consumeWhitespace()
	}

	// end of input
	tok := newTokenByte(token.EOF, l.ch)
	l.tokens <- tok
	close(l.tokens)
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

// Generates a string from a consecutive sequence of characters where isAlpha
// returns true starting with the current character (l.ch) and continuing until
// the peekCharacter does not meet the isAlpha condition. When returning, l.ch
// will point to the last consumed character.
func (l *Lexer) readWord() string {
	sb := strings.Builder{}

	sb.WriteByte(l.ch)
	for isAlpha(l.peekCharacter()) {
		l.readCharacter()
		sb.WriteByte(l.ch)
	}

	return sb.String()
}

// Generates a string from a consecutive sequence of characters where isDigit
// returns true starting with the current character (l.ch) and continuing until
// the peekCharacter does not meet the isDigit condition. When returning, l.ch
// will point to the last consumed character.
func (l *Lexer) readNumber() string {
	sb := strings.Builder{}

	sb.WriteByte(l.ch)
	for isDigit(l.peekCharacter()) {
		l.readCharacter()
		sb.WriteByte(l.ch)
	}

	return sb.String()
}
