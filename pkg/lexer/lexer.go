package lexer

import (
	"strings"

	"github.com/freddiehaddad/corrosion/pkg/token"
)

// The buffer size of the Lexer's token channel.
const capacity = 10

// ----------------------------------------------------------------------------
// Lexer
// ----------------------------------------------------------------------------

// The Lexer object represents the state of the lexer.
type Lexer struct {
	tokens       chan token.Token // tokens generated from stream
	input        string           // input stream
	position     int              // the current position in the stream
	readPosition int              // the next position to read from
	ch           byte             // the character at position
}

// Creates and returns a Lexer.
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.tokens = make(chan token.Token, capacity)
	l.readCharacter()
	go l.generateTokens()
	return l
}

// ----------------------------------------------------------------------------
// Token generators
// ----------------------------------------------------------------------------

// Create a new token for single-character terminals.
func newTokenByte(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// Create a new token for string literals.
func newTokenString(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
	}
}

// Returns the next token from the buffered channel.
func (l *Lexer) NextToken() token.Token {
	tok := <-l.tokens
	return tok
}

// Loop through the input stream pushing tokens to the buffered tokens channel.
func (l *Lexer) generateTokens() {
	l.consumeWhitespace()

	for l.ch != token.EOF_VALUE {
		var tok token.Token

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
		case '{':
			tok = newTokenByte(token.LBRACE, l.ch)
		case '}':
			tok = newTokenByte(token.RBRACE, l.ch)

		// operators
		case '-':
			tok = newTokenByte(token.MINUS, l.ch)
		case '+':
			tok = newTokenByte(token.PLUS, l.ch)
		case '*':
			tok = newTokenByte(token.MULTIPLY, l.ch)
		case '/':
			tok = newTokenByte(token.DIVIDE, l.ch)
		case '=':
			if l.peekCharacter() == '=' {
				l.readCharacter()
				tok = newTokenString(token.EQ, "==")
			} else {
				tok = newTokenByte(token.ASSIGN, l.ch)
			}
		case '!':
			if l.peekCharacter() == '=' {
				l.readCharacter()
				tok = newTokenString(token.NOT_EQ, "!=")
			} else {
				tok = newTokenByte(token.BANG, l.ch)
			}
		case '<':
			if l.peekCharacter() == '=' {
				l.readCharacter()
				tok = newTokenString(token.LT_EQUAL, "<=")
			} else {
				tok = newTokenByte(token.LT, l.ch)
			}
		case '>':
			if l.peekCharacter() == '=' {
				l.readCharacter()
				tok = newTokenString(token.GT_EQUAL, ">=")
			} else {
				tok = newTokenByte(token.GT, l.ch)
			}

		default:
			// identifiers
			if isAlpha(l.ch) {
				s := l.readWord()
				tt := token.LookupType(s)
				tok = newTokenString(tt, s)
				// integer literals
			} else if isDigit(l.ch) {
				s := l.readNumber()
				tok = newTokenString(token.INTEGER, s)
				// invalid tokens
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

// ----------------------------------------------------------------------------
// Helper functions for token generation
// ----------------------------------------------------------------------------

// Returns true if ch is a numeric value.
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// Returns true if ch is an alpha value.
func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z'
}

// Returns true if ch is a whitespace value.
func isWhitespace(ch byte) bool {
	switch ch {
	case ' ', '\r', '\n', '\t':
		return true
	default:
		return false
	}
}

// Consume all sequential whitespace characters in the lexer input stream.
func (l *Lexer) consumeWhitespace() {
	for isWhitespace(l.ch) {
		l.readCharacter()
	}
}

// Advances the lexer position one character.  If the lexer has reached the end
// of the stream, no change to the state occurs.
func (l *Lexer) readCharacter() {
	if l.eof() {
		l.ch = token.EOF_VALUE
		return
	}

	l.ch = l.input[l.readPosition]
	l.position = l.readPosition
	l.readPosition++
}

// Returns the next character in the sequence without advancing.  Returns
// the end of file value if the stream has reached the end.
func (l *Lexer) peekCharacter() byte {
	if l.eof() {
		return token.EOF_VALUE
	}

	return l.input[l.readPosition]
}

// Returns true if the lexer has reached the end of input.
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
