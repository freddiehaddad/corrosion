package evaluator

import (
	"testing"

	"github.com/freddiehaddad/corrosion/pkg/lexer"
	"github.com/freddiehaddad/corrosion/pkg/object"
	"github.com/freddiehaddad/corrosion/pkg/parser"
)

func TestEvalIntegerExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5;", "5"},
		{"10;", "10"},
		{"-10;", "-10"},
		{"--10;", "10"},
	}

	e := object.NewEnvironment()

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Integer:
			testIntegerObject(t, obj, test.expected)
		default:
			t.Errorf("object is not Integer. got=%T (%+v)",
				obj, obj)
		}
	}
}

func TestArithmeticExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{" 1 +  8;", "9"},
		{"-2 +  6;", "4"},
		{" 3 + -6;", "-3"},
		{"-4 + -5;", "-9"},
		{" 5 -  4;", "1"},
		{"-6 -  3;", "-9"},
		{" 7 - -2;", "9"},
		{"-8 - -1;", "-7"},
	}

	e := object.NewEnvironment()

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Integer:
			testIntegerObject(t, obj, test.expected)
		default:
			t.Errorf("object is not Integer. got=%T (%+v)",
				obj, obj)
		}
	}
}

func TestVariableDeclaration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"int x = 3; x;", "3"},
		{"int y = 2; y;", "2"},
		{"x + y;", "5"},
		{"y + x;", "5"},
		{"x - y;", "1"},
		{"y - x;", "-1"},
	}

	e := object.NewEnvironment()

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Integer:
			testIntegerObject(t, obj, test.expected)
		default:
			t.Errorf("object is not Integer. got=%T (%+v)",
				obj, obj)
		}
	}
}

func testIntegerObject(t *testing.T, obj *object.Integer, expected string) {
	if obj.Value != expected {
		t.Errorf("object has wrong value. got=%s, expected=%s",
			obj.Value, expected)
	}
}
