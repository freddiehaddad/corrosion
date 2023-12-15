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

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Integer:
			testIntegerObject(t, index, obj, test.expected)
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

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Boolean:
			testBooleanObject(t, index, obj, test.expected)
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

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Boolean:
			testBooleanObject(t, index, obj, test.expected)
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

		{"true != true;", false},
		{"true != false;", true},
		{"false != true;", true},
		{"false != false;", false},
		{"!true != false;", false},
		{"true != !false;", false},
		{"!true != !false;", true},

		{"5 != 5;", false},
		{"2 != 5;", true},
		{"2 + 3 != 5;", false},
		{"8 != 10 - 2;", false},
		{"2 - 3 != 5;", true},
		{"7 != 10 - 2;", true},

		{"5 < 5;", false},
		{"2 < 5;", true},
		{"2 + 3 < 5;", false},
		{"8 < 10 - 2;", false},
		{"2 - 3 < 5;", true},
		{"7 < 10 - 2;", true},

		{"5 <= 5;", true},
		{"2 <= 5;", true},
		{"2 + 3 <= 5;", true},
		{"8 <= 10 - 2;", true},
		{"2 - 3 <= 5;", true},
		{"7 <= 10 - 2;", true},

		{"5 > 5;", false},
		{"2 > 5;", false},
		{"2 + 3 > 5;", false},
		{"8 > 10 - 2;", false},
		{"2 - 3 > 5;", false},
		{"7 > 10 - 2;", false},

		{"5 >= 5;", true},
		{"2 >= 5;", false},
		{"2 + 3 >= 5;", true},
		{"8 >= 10 - 2;", true},
		{"2 - 3 >= 5;", false},
		{"7 >= 10 - 2;", false},

		{"2 + 4 < 3 + 4 != false;", true},
		{"2 * 4 < 3 * 4 != false;", true},
		{"2 + 4 <= 3 + 4 != false;", true},
		{"2 * 4 <= 3 * 4 != false;", true},

		{"2 + 4 > 3 + 4 != true;", true},
		{"2 * 4 > 3 * 4 != true;", true},
		{"2 + 4 >= 3 + 4 != true;", true},
		{"2 * 4 >= 3 * 4 != true;", true},
	}

	e := object.NewEnvironment()

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Boolean:
			testBooleanObject(t, index, obj, test.expected)
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

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Integer:
			testIntegerObject(t, index, obj, test.expected)
		default:
			t.Errorf("object is not Integer. got=%T (%+v)",
				obj, obj)
		}
	}
}

func TestGroupedExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"(2 + 3)", 5},
		{"-(-2 + -3)", 5},
		{"-(-2 - -5)", -3},
		{"-(2 + 3);", -5},
		{"(2 + 3) * 4;", 20},
		{"2 * (3 + 4);", 14},
		{"(201 - 1) / (2 * (2 + 3));", 20},

		{"!(true == false);", true},
	}

	e := object.NewEnvironment()

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, index, result, int64(expected))
		case bool:
			testBooleanObject(t, index, result, expected)
		default:
			t.Errorf("test[%d]: unsupported type for expected=%T (%+v)", index, test.expected, test.expected)
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

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Integer:
			testIntegerObject(t, index, obj, test.expected)
		default:
			t.Errorf("object is not Integer. got=%T (%+v)",
				obj, obj)
		}
	}
}

func TestAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"int x = 4; x = 5; x;", 5},
		{"int foo = 4 + 3; foo = foo * 2; foo;", 14},
		{"int foo = 10; foo = foo * foo; foo;", 100},
	}

	for index, test := range tests {
		l := lexer.New(test.input)
		p := parser.New(l)
		program := p.ParseProgram()
		e := object.NewEnvironment()

		result := Eval(program, e)

		switch obj := result.(type) {
		case *object.Integer:
			testIntegerObject(t, index, obj, test.expected)
		default:
			t.Errorf("object is not Integer. got=%T (%+v)",
				obj, obj)
		}
	}
}

func testBooleanObject(t *testing.T, index int, obj object.Object, expected bool) {
	switch o := obj.(type) {
	case *object.Boolean:
		if o.Value != expected {
			t.Errorf("tests[%d]: object has wrong value. got=%t, expected=%t",
				index, o.Value, expected)
		}
	default:
		t.Errorf("tests[%d]: wrong type. got=%T (%+v)", index, obj, obj)
	}
}

func testIntegerObject(t *testing.T, index int, obj object.Object, expected int64) {
	switch o := obj.(type) {
	case *object.Integer:
		if o.Value != expected {
			t.Errorf("tests[%d]: object has wrong value. got=%d, expected=%d",
				index, o.Value, expected)
		}
	default:
		t.Errorf("tests[%d]: wrong type. got=%T (%+v)", index, obj, obj)
	}
}

func testErrorObject(t *testing.T, obj *object.Error, expected string) {
	if obj.Value != expected {
		t.Errorf("object has wrong value. got=%s, expected=%s",
			obj.Value, expected)
	}
}
