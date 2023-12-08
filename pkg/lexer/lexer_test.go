package lexer

import (
	"testing"

	"github.com/freddiehaddad/corrosion/pkg/token"
)

type expectedToken struct {
	expectedType    token.TokenType
	expectedLiteral string
}

func compareToken(t *testing.T, i int, et expectedToken, nt token.Token) {
	if nt.Type != et.expectedType {
		t.Errorf("tests[%d] - token type wrong: expected=%q got=%q",
			i, et.expectedType, nt.Type)
	}
	if nt.Literal != et.expectedLiteral {
		t.Errorf("tests[%d] - token literal wrong: expected=%q got=%q",
			i, et.expectedLiteral, nt.Literal)
	}
}

func compareTokens(t *testing.T, l *Lexer, tokens []expectedToken) {
	for i, tt := range tokens {
		tok := l.NextToken()
		compareToken(t, i, tt, tok)
	}
}

func TestNextToken(t *testing.T) {
	input := "int return x +-*= 10;$"
	tests := []expectedToken{
		{expectedType: token.INT, expectedLiteral: "int"},
		{expectedType: token.RETURN, expectedLiteral: "return"},
		{expectedType: token.IDENT, expectedLiteral: "x"},
		{expectedType: token.PLUS, expectedLiteral: "+"},
		{expectedType: token.MINUS, expectedLiteral: "-"},
		{expectedType: token.MULTIPLY, expectedLiteral: "*"},
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.INTEGER, expectedLiteral: "10"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.ILLEGAL, expectedLiteral: "$"},
		{
			expectedType:    token.EOF,
			expectedLiteral: string(token.EOF_VALUE),
		},
	}

	l := New(input)
	compareTokens(t, l, tests)
}
