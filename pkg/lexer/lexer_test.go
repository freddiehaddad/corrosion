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
	input := "var return if else x true false !!= <<= >>= +-*/= 10;==)({},$"
	tests := []expectedToken{
		{expectedType: token.VAR, expectedLiteral: "var"},
		{expectedType: token.RETURN, expectedLiteral: "return"},
		{expectedType: token.IF, expectedLiteral: "if"},
		{expectedType: token.ELSE, expectedLiteral: "else"},
		{expectedType: token.IDENT, expectedLiteral: "x"},
		{expectedType: token.TRUE, expectedLiteral: "true"},
		{expectedType: token.FALSE, expectedLiteral: "false"},
		{expectedType: token.BANG, expectedLiteral: "!"},
		{expectedType: token.NOT_EQ, expectedLiteral: "!="},
		{expectedType: token.LT, expectedLiteral: "<"},
		{expectedType: token.LT_EQUAL, expectedLiteral: "<="},
		{expectedType: token.GT, expectedLiteral: ">"},
		{expectedType: token.GT_EQUAL, expectedLiteral: ">="},
		{expectedType: token.PLUS, expectedLiteral: "+"},
		{expectedType: token.MINUS, expectedLiteral: "-"},
		{expectedType: token.MULTIPLY, expectedLiteral: "*"},
		{expectedType: token.DIVIDE, expectedLiteral: "/"},
		{expectedType: token.ASSIGN, expectedLiteral: "="},
		{expectedType: token.INTEGER, expectedLiteral: "10"},
		{expectedType: token.SEMICOLON, expectedLiteral: ";"},
		{expectedType: token.EQ, expectedLiteral: "=="},
		{expectedType: token.RPAREN, expectedLiteral: ")"},
		{expectedType: token.LPAREN, expectedLiteral: "("},
		{expectedType: token.LBRACE, expectedLiteral: "{"},
		{expectedType: token.RBRACE, expectedLiteral: "}"},
		{expectedType: token.COMMA, expectedLiteral: ","},
		{expectedType: token.ILLEGAL, expectedLiteral: "$"},
		{
			expectedType:    token.EOF,
			expectedLiteral: string(token.EOF_VALUE),
		},
	}

	l := New(input)
	compareTokens(t, l, tests)
}
