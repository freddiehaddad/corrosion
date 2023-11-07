package lexer

import "github.com/freddiehaddad/corrosion/pkg/token"

type Lexer struct {
	input    string
	position int
	ch       byte
}

func New(input string) *Lexer {
	return &Lexer{input: input, position: 0, ch: 0}
}

func (l *Lexer) NextToken() token.Token {
	tok := token.Token{}

	l.readCharacter()
	switch l.ch {
	// delimiters
	case ';':
		tok.Type = token.SEMICOLON
		tok.Literal = token.SEMICOLON
	case ',':
		tok.Type = token.COMMA
		tok.Literal = token.COMMA
	case '(':
		tok.Type = token.LPAREN
		tok.Literal = token.LPAREN
	case ')':
		tok.Type = token.RPAREN
		tok.Literal = token.RPAREN
	case '[':
		tok.Type = token.LBRACKET
		tok.Literal = token.LBRACKET
	case ']':
		tok.Type = token.RBRACKET
		tok.Literal = token.RBRACKET
	case '{':
		tok.Type = token.LBRACE
		tok.Literal = token.LBRACE
	case '}':
		tok.Type = token.RBRACE
		tok.Literal = token.RBRACE
	// operators
	case '=':
		tok.Type = token.ASSIGN
		tok.Literal = token.ASSIGN
	case '+':
		tok.Type = token.PLUS
		tok.Literal = token.PLUS
	case '-':
		tok.Type = token.MINUS
		tok.Literal = token.MINUS
	case '*':
		tok.Type = token.ASTERISK
		tok.Literal = token.ASTERISK
	case '/':
		tok.Type = token.FSLASH
		tok.Literal = token.FSLASH
	// end of input
	case 0:
		tok.Type = token.EOF
		tok.Literal = token.EOF
	// errors
	default:
		tok.Type = token.UNKNOWN
		tok.Literal = string(l.ch)
	}

	return tok
}

func (l *Lexer) readCharacter() {
	if l.eof() {
		l.ch = 0
		return
	}

	l.ch = l.input[l.position]
	l.position++
}

func (l *Lexer) eof() bool {
	return l.position == len(l.input)
}
