package lexer

import (
	"testing"

	"github.com/freddiehaddad/corrosion/pkg/token"
)

func compareToken(t *testing.T, i int, et, nt token.Token) {
	if nt.Type != et.Type {
		t.Errorf("tests[%d] - token type wrong: expected=%q got=%q", i, et.Type, nt.Type)
	}
	if nt.Literal != et.Literal {
		t.Errorf("tests[%d] - token literal wrong: expected=%q got=%q", i, et.Literal, nt.Literal)
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
		{`;`, true},
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
		{"", []byte{token.EOF_VALUE}},
		{";", []byte{token.EOF_VALUE}},
		{"()", []byte{')', token.EOF_VALUE}},
	}

	for index, test := range tests {
		l := New(test.input)
		for expectedCharacterIndex, expectedCharacter := range test.expectedCharacters {
			result := l.peekCharacter()
			l.readCharacter()

			if result != expectedCharacter {
				t.Errorf("tests[%d] - peek character from position=%d wrong. expected=%q got=%q\n",
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
		{"", []byte{token.EOF_VALUE}},
		{"a", []byte{'a', token.EOF_VALUE}},
		{"ab", []byte{'a', 'b', token.EOF_VALUE}},
	}

	for index, test := range tests {
		l := New(test.input)
		for expectedCharacterIndex, expectedCharacter := range test.expectedCharacters {
			if l.ch != expectedCharacter {
				t.Errorf("tests[%d] - character at position=%d wrong. expected=%q got=%q\n",
					index, expectedCharacterIndex, expectedCharacter, l.ch)
			}
			l.readCharacter()
		}
	}
}

func TestReadCharacterPositions(t *testing.T) {
	tests := []struct {
		input     string
		positions []int
	}{
		// values are in pairs (position,  readPosition)
		{"", []int{0}},
		{"a", []int{0, 1, 0, 1}},
		{"ab", []int{0, 1, 1, 2}},
	}

	for index, test := range tests {
		l := New(test.input)
		for position := 0; position < len(test.positions)-1; position += 2 {
			expectedPosition := test.positions[position]
			expectedReadPosition := test.positions[position+1]

			if l.position != expectedPosition {
				t.Errorf("tests[%d] - position wrong. expected=%d got=%d\n", index, expectedPosition, l.position)
			}

			if l.readPosition != expectedReadPosition {
				t.Errorf("tests[%d] - readPosition wrong. expected=%d got=%d\n", index, expectedReadPosition, l.readPosition)
			}

			l.readCharacter()
		}
	}
}

func TestIsWhitespace(t *testing.T) {
	input := " \r\n\t"
	tests := []bool{true, true, true, true, true, true}

	for index := 0; index < len(input); index++ {
		ch := input[index]
		expected := tests[index]
		result := isWhitespace(ch)
		if expected != result {
			t.Errorf("tests[%d] - whitespace check wrong for %q. expected=%t got=%t\n", index, ch, expected, result)
		}
	}
}

func TestIgnoreWhitespace(t *testing.T) {
	input := ",\r\n\t ;"
	expectedPositions := []int{0, 5}

	l := New(input)
	for index, expectedPosition := range expectedPositions {
		l.skipWhitespace()
		if l.position != expectedPosition {
			t.Errorf("tests[%d] - position wrong. expected=%d, got=%d\n", index, expectedPosition, l.position)
		}
		if l.ch != input[l.position] {
			t.Errorf("tests[%d] - character wrong. expected=%c, got=%c\n", index, input[l.position], l.position)
		}
		l.readCharacter()
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

func TestTwoDigitOperatorTokens(t *testing.T) {
	input := `==!=&&||`
	tests := []token.Token{
		{Type: token.EQUAL, Literal: "=="},
		{Type: token.NOT_EQUAL, Literal: "!="},
		{Type: token.AND, Literal: "&&"},
		{Type: token.OR, Literal: "||"},
	}

	l := New(input)
	compareTokens(t, l, tests)
}

func TestSingleCharaterTokens(t *testing.T) {
	input := `;,()[]{}=+-*/<>&|^~!`
	tests := []token.Token{
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACKET, Literal: "["},
		{Type: token.RBRACKET, Literal: "]"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.FSLASH, Literal: "/"},
		{Type: token.LESS_THAN, Literal: "<"},
		{Type: token.GREATER_THAN, Literal: ">"},
		{Type: token.BITWISE_AND, Literal: "&"},
		{Type: token.BITWISE_OR, Literal: "|"},
		{Type: token.BITWISE_XOR, Literal: "^"},
		{Type: token.BITWISE_NOT, Literal: "~"},
		{Type: token.BANG, Literal: "!"},
		{Type: token.EOF, Literal: string(token.EOF_VALUE)},
	}
	l := New(input)

	compareTokens(t, l, tests)
}

func TestLeadingWhitespace(t *testing.T) {
	input := " ;"
	tests := []token.Token{
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: string(token.EOF_VALUE)},
	}
	l := New(input)

	compareTokens(t, l, tests)
}

func TestTrailingWhitespace(t *testing.T) {
	input := "; "
	tests := []token.Token{
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: string(token.EOF_VALUE)},
	}
	l := New(input)

	compareTokens(t, l, tests)
}

func TestMixedWhitespace(t *testing.T) {
	input := "; , ( ) [ ] { } = + - * / < > & | ^ ~ !"
	tests := []token.Token{
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACKET, Literal: "["},
		{Type: token.RBRACKET, Literal: "]"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.FSLASH, Literal: "/"},
		{Type: token.LESS_THAN, Literal: "<"},
		{Type: token.GREATER_THAN, Literal: ">"},
		{Type: token.BITWISE_AND, Literal: "&"},
		{Type: token.BITWISE_OR, Literal: "|"},
		{Type: token.BITWISE_XOR, Literal: "^"},
		{Type: token.BITWISE_NOT, Literal: "~"},
		{Type: token.BANG, Literal: "!"},
		{Type: token.EOF, Literal: string(token.EOF_VALUE)},
	}
	l := New(input)

	compareTokens(t, l, tests)
}
