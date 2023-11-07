package lexer

import (
	"testing"

	"github.com/freddiehaddad/corrosion/pkg/token"
)

func compareToken(t *testing.T, i int, et, nt token.Token) {
	if nt.Type != et.Type {
		t.Fatalf("tests[%d] - token type wrong: expected=%q got=%q", i, et.Type, nt.Type)
	}
	if nt.Literal != et.Literal {
		t.Fatalf("tests[%d] - token literal wrong: expected=%q got=%q", i, et.Literal, nt.Literal)
	}
}

func compareTokens(t *testing.T, l *Lexer, tokens []token.Token) {
	for i, tt := range tokens {
		tok := l.NextToken()
		compareToken(t, i, tt, tok)
	}
}

func TestEmptyInput(t *testing.T) {
	input := ``
	tests := []token.Token{
		{Type: token.EOF, Literal: token.EOF},
	}
	l := New(input)

	compareTokens(t, l, tests)
}

func TestUnknownInput(t *testing.T) {
	input := `@`
	tests := []token.Token{
		{Type: token.UNKNOWN, Literal: "@"},
		{Type: token.EOF, Literal: token.EOF},
	}
	l := New(input)

	compareTokens(t, l, tests)

}

func TestSingleCharaterTokens(t *testing.T) {
	input := `;,()[]{}=+-*/`
	tests := []token.Token{
		{Type: token.SEMICOLON, Literal: token.SEMICOLON},
		{Type: token.COMMA, Literal: token.COMMA},
		{Type: token.LPAREN, Literal: token.LPAREN},
		{Type: token.RPAREN, Literal: token.RPAREN},
		{Type: token.LBRACKET, Literal: token.LBRACKET},
		{Type: token.RBRACKET, Literal: token.RBRACKET},
		{Type: token.LBRACE, Literal: token.LBRACE},
		{Type: token.RBRACE, Literal: token.RBRACE},
		{Type: token.ASSIGN, Literal: token.ASSIGN},
		{Type: token.PLUS, Literal: token.PLUS},
		{Type: token.MINUS, Literal: token.MINUS},
		{Type: token.ASTERISK, Literal: token.ASTERISK},
		{Type: token.FSLASH, Literal: token.FSLASH},
		{Type: token.EOF, Literal: token.EOF},
	}
	l := New(input)

	compareTokens(t, l, tests)
}
