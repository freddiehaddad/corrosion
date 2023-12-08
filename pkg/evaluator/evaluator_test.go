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

		{" 5 *  4;", "20"},
		{"-6 *  3;", "-18"},
		{" 7 * -2;", "-14"},
		{"-8 * -1;", "8"},

		{" 5 /  4;", "1"},
		{"-6 /  3;", "-2"},
		{" 7 / -2;", "-3"},
		{"-8 / -1;", "8"},

		{" 5 * 10 + 15;", "65"},
		{" 5 * 10 - 15;", "35"},
		{" 5 * 10 + -15;", "35"},
		{" 5 * 10 - -15;", "65"},

		{" -5 * 10 + 15;", "-35"},
		{" -5 * 10 - 15;", "-65"},
		{" -5 * 10 + -15;", "-65"},
		{" -5 * 10 - -15;", "-35"},

		{" 5 * -10 + 15;", "-35"},
		{" 5 * -10 - 15;", "-65"},
		{" 5 * -10 + -15;", "-65"},
		{" 5 * -10 - -15;", "-35"},

		{" -5 * -10 - 15;", "35"},
		{" -5 * -10 + -15;", "35"},

		{" 10 / 2 * 7;", "35"},
		{" 10 * 2 / 4;", "5"},
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

func TestDivideByZeroError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0 / 0;", "ERROR: divide by zero error in expression ((0 / 0))"},
		{"1 / 0;", "ERROR: divide by zero error in expression ((1 / 0))"},
		{"1 + 2 / 0;", "ERROR: divide by zero error in expression ((2 / 0))"},
		{"1 / 0 + 2;", "ERROR: divide by zero error in expression ((1 / 0))"},
	}

	e := object.NewEnvironment()

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Error:
			testErrorObject(t, obj, test.expected)
		default:
			t.Errorf("object is not Error. got=%T (%+v)",
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

func testErrorObject(t *testing.T, obj *object.Error, expected string) {
	if obj.Value != expected {
		t.Errorf("object has wrong value. got=%s, expected=%s",
			obj.Value, expected)
	}
}
