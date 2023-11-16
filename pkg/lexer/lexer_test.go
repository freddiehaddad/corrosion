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

func TestEof(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{``, true},
		{` `, false},
	}

	for index, test := range tests {
		l := New(test.input)
		result := l.eof()
		if result != test.expected {
			t.Errorf("tests[%d] - eof wrong: expected=%v got=%v\n", index, test.expected, result)
		}
	}
}

func TestPeekCharacter(t *testing.T) {
	tests := []struct {
		input              string
		expectedCharacters []uint8
	}{
		{``, []byte{token.EOF_VALUE}},
		{`;`, []byte{';', token.EOF_VALUE}},
		{`()`, []byte{'(', ')', token.EOF_VALUE}},
	}

	for index, test := range tests {
		l := New(test.input)
		for expectedCharacterIndex, expectedCharacter := range test.expectedCharacters {
			result := l.peekCharacter()
			l.readCharacter()

			if result != expectedCharacter {
				t.Errorf("tests[%d] - peek character from position=%d wrong. expected=%b got=%b\n",
					index, expectedCharacterIndex, expectedCharacter, result)
			}
		}
	}
}

func TestReadCharacter(t *testing.T) {
	tests := []struct {
		input              string
		expectedCharacters []uint8
	}{
		{``, []byte{token.EOF_VALUE}},
		{`a`, []byte{'a', token.EOF_VALUE}},
		{`ab`, []byte{'a', 'b', token.EOF_VALUE}},
	}

	for index, test := range tests {
		l := New(test.input)
		for expectedCharacterIndex, expectedCharacter := range test.expectedCharacters {
			l.readCharacter()

			if l.ch != expectedCharacter {
				t.Errorf("tests[%d] - character at position=%d wrong. expected=%b got=%b\n",
					index, expectedCharacterIndex, expectedCharacter, l.ch)
			}
		}
	}
}

func TestReadCharacterPositions(t *testing.T) {
	tests := []struct {
		input     string
		positions []int
	}{
		// values are in pairs (position,  readPosition)
		{``, []int{0}},
		{`a`, []int{0, 1, 0, 1}},
		{`ab`, []int{0, 1, 1, 2}},
	}

	for index, test := range tests {
		l := New(test.input)
		for position := 0; position < len(test.positions)-1; position += 2 {
			l.readCharacter()

			expectedPosition := test.positions[position]
			expectedReadPosition := test.positions[position+1]

			if l.position != expectedPosition {
				t.Errorf("tests[%d] - position wrong. expected=%d got=%d\n", index, expectedPosition, l.position)
			}

			if l.readPosition != expectedReadPosition {
				t.Errorf("tests[%d] - readPosition wrong. expected=%d got=%d\n", index, expectedReadPosition, l.readPosition)
			}
		}
	}
}

func TestEmptyInput(t *testing.T) {
	input := ``
	tests := []token.Token{
		{Type: token.EOF, Literal: string(token.EOF_VALUE)},
	}
	l := New(input)

	compareTokens(t, l, tests)
}

func TestUnknownInput(t *testing.T) {
	input := `@`
	tests := []token.Token{
		{Type: token.UNKNOWN, Literal: "@"},
		{Type: token.EOF, Literal: string(token.EOF_VALUE)},
	}
	l := New(input)

	compareTokens(t, l, tests)

}

func TestMultiDigitTokens(t *testing.T) {
	input := `==!=&&||`
	tests := []token.Token{
		{Type: token.EQUAL, Literal: token.EQUAL},
		{Type: token.NOT_EQUAL, Literal: token.NOT_EQUAL},
		{Type: token.AND, Literal: token.AND},
		{Type: token.OR, Literal: token.OR},
	}

	l := New(input)
	compareTokens(t, l, tests)
}

func TestSingleCharaterTokens(t *testing.T) {
	input := `;,()[]{}=+-*/<>&|^~!`
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
		{Type: token.LESS_THAN, Literal: token.LESS_THAN},
		{Type: token.GREATER_THAN, Literal: token.GREATER_THAN},
		{Type: token.BITWISE_AND, Literal: token.BITWISE_AND},
		{Type: token.BITWISE_OR, Literal: token.BITWISE_OR},
		{Type: token.BITWISE_XOR, Literal: token.BITWISE_XOR},
		{Type: token.BITWISE_NOT, Literal: token.BITWISE_NOT},
		{Type: token.BANG, Literal: token.BANG},
		{Type: token.EOF, Literal: string(token.EOF_VALUE)},
	}
	l := New(input)

	compareTokens(t, l, tests)
}
