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
		expected int64
	}{
		{"5;", 5},
		{"10;", 10},
		{"-10;", -10},
		{"--10;", 10},
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

func TestEvalNotExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true;", false},
		{"!false;", true},
		{"!!true;", true},
		{"!!false;", false},
	}

	e := object.NewEnvironment()

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Boolean:
			testBooleanObject(t, obj, test.expected)
		default:
			t.Errorf("object is not Boolean. got=%T (%+v)",
				obj, obj)
		}
	}
}

func TestEvalBooleanExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	e := object.NewEnvironment()

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Boolean:
			testBooleanObject(t, obj, test.expected)
		default:
			t.Errorf("object is not Boolean. got=%T (%+v)",
				obj, obj)
		}
	}
}

func TestComparisonExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true == true;", true},
		{"true == false;", false},
		{"false == true;", false},
		{"false == false;", true},

		{"!true == false;", true},
		{"true == !false;", true},
		{"!true == !false;", false},

		{"5 == 5;", true},
		{"2 == 5;", false},

		{"2 + 3 == 5;", true},
		{"8 == 10 - 2;", true},

		{"2 - 3 == 5;", false},
		{"7 == 10 - 2;", false},
	}

	e := object.NewEnvironment()

	for _, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Boolean:
			testBooleanObject(t, obj, test.expected)
		default:
			t.Errorf("object is not Boolean. got=%T (%+v)",
				obj, obj)
		}
	}
}

func TestArithmeticExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{" 1 +  8;", 9},
		{"-2 +  6;", 4},
		{" 3 + -6;", -3},
		{"-4 + -5;", -9},

		{" 5 -  4;", 1},
		{"-6 -  3;", -9},
		{" 7 - -2;", 9},
		{"-8 - -1;", -7},

		{" 5 *  4;", 20},
		{"-6 *  3;", -18},
		{" 7 * -2;", -14},
		{"-8 * -1;", 8},

		{" 5 /  4;", 1},
		{"-6 /  3;", -2},
		{" 7 / -2;", -3},
		{"-8 / -1;", 8},

		{" 5 * 10 + 15;", 65},
		{" 5 * 10 - 15;", 35},
		{" 5 * 10 + -15;", 35},
		{" 5 * 10 - -15;", 65},

		{" -5 * 10 + 15;", -35},
		{" -5 * 10 - 15;", -65},
		{" -5 * 10 + -15;", -65},
		{" -5 * 10 - -15;", -35},

		{" 5 * -10 + 15;", -35},
		{" 5 * -10 - 15;", -65},
		{" 5 * -10 + -15;", -65},
		{" 5 * -10 - -15;", -35},

		{" -5 * -10 - 15;", 35},
		{" -5 * -10 + -15;", 35},

		{" 10 / 2 * 7;", 35},
		{" 10 * 2 / 4;", 5},
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

func TestArithmeticExpressionsWithVariables(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"3+x;", "ERROR: undefined identifier=\"x\" (x)"},
		{"x+3;", "ERROR: undefined identifier=\"x\" (x)"},
		{"int x = 3; x+y;", "ERROR: undefined identifier=\"y\" (y)"},
		{"int y = 3; x+y;", "ERROR: undefined identifier=\"x\" (x)"},
	}

	for _, test := range tests {
		e := object.NewEnvironment()

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

func TestDivideByZeroError(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"0 / 0;",
			"ERROR: divide by zero error in expression (0 / 0)",
		},
		{
			"1 / 0;",
			"ERROR: divide by zero error in expression (1 / 0)",
		},
		{
			"1 + 2 / 0;",
			"ERROR: divide by zero error in expression (2 / 0)",
		},
		{
			"1 / 0 + 2;",
			"ERROR: divide by zero error in expression (1 / 0)",
		},
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
		expected int64
	}{
		{"int x = 3; x;", 3},
		{"int y = 2; y;", 2},
		{"x + y;", 5},
		{"y + x;", 5},
		{"x - y;", 1},
		{"y - x;", -1},
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

func testBooleanObject(t *testing.T, obj *object.Boolean, expected bool) {
	if obj.Value != expected {
		t.Errorf("object has wrong value. got=%t, expected=%t",
			obj.Value, expected)
	}
}

func testIntegerObject(t *testing.T, obj *object.Integer, expected int64) {
	if obj.Value != expected {
		t.Errorf("object has wrong value. got=%d, expected=%d",
			obj.Value, expected)
	}
}

func testErrorObject(t *testing.T, obj *object.Error, expected string) {
	if obj.Value != expected {
		t.Errorf("object has wrong value. got=%s, expected=%s",
			obj.Value, expected)
	}
}
