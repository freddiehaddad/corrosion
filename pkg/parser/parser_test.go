package parser

import (
	"strconv"
	"testing"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/lexer"
)

type testResults [][]string

// ----------------------------------------------------------------------------
// Test helper functions
// ----------------------------------------------------------------------------

func checkProgram(t *testing.T, p *ast.Program) {
	if p == nil {
		t.Fatalf("ParseProgram returned nil\n")
	}
}

func checkErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("ParseProgram returned %d errors\n", len(errors))
	for errno, error := range errors {
		t.Errorf("errors[%d]: %s", errno, error)
	}

	t.FailNow()
}

func checkLength(t *testing.T, length int, stmts []ast.Statement) {
	if len(stmts) != length {
		t.Fatalf("ParseProgram returned %d statements, expected %d\n",
			len(stmts), length)
	}
}

func checkStatements(t *testing.T, expects testResults, stmts []ast.Statement) {
	for i, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.DeclarationStatement:
			checkDeclarationStatement(t, i, expects[i], s)
		case *ast.ReturnStatement:
			checkReturnStatement(t, i, expects[i], s)
		case *ast.ExpressionStatement:
			checkExpressionStatement(t, i, expects[i], s)
		default:
			t.Errorf("tests[%d]: wrong type. got=%T\n", i, stmt)
		}
	}
}

func checkBoolean(
	t *testing.T, test int, expected []string, node *ast.Boolean,
) {
	if node.Token.Literal != expected[0] {
		t.Errorf("tests[%d]: incorrect value. expected=%q got=%q\n",
			test, expected[0], node.Token.Literal)
	}

	value, err := strconv.ParseBool(expected[0])
	if err != nil {
		t.Errorf("tests[%d]: failed to parse expected[0]=%s: err=%s",
			test, expected[0], err.Error())
	}

	if node.Value != value {
		t.Errorf("tests[%d]: incorrect value. expected=%t got=%t\n",
			test, value, node.Value)
	}
}

func checkIdentifier(
	t *testing.T, test int, expected []string, node *ast.Identifier,
) {
	if node.Value != expected[0] {
		t.Errorf("tests[%d]: incorrect value. expected=%q got=%q\n",
			test, expected[0], node.Value)
	}
}

func checkIntegerLiteral(
	t *testing.T, test int, expected []string, node *ast.IntegerLiteral,
) {
	if node.Token.Literal != expected[0] {
		t.Errorf("tests[%d]: incorrect value. expected=%q got=%q\n",
			test, expected[0], node.Value)
	}

	value, err := strconv.ParseInt(expected[0], 0, 64)
	if err != nil {
		t.Errorf("tests[%d]: failed to parse expected[0]=%s: err=%s",
			test, expected[0], err.Error())
	}

	if node.Value != value {
		t.Errorf("tests[%d]: incorrect value. expected=%d got=%d\n",
			test, value, node.Value)
	}
}

func checkInfixExpression(
	t *testing.T, test int, expected []string, node *ast.InfixExpression,
) {
	if node.Left.String() != expected[0] {
		t.Errorf("tests[%d]: incorrect Expression. expected=%q got=%q",
			test, expected[0], node.Left.String())
	}

	if node.Operator != expected[1] {
		t.Errorf("tests[%d]: incorrect operator. expected=%q got=%q",
			test, expected[1], node.Operator)
	}

	if node.Right.String() != expected[2] {
		t.Errorf("tests[%d]: incorrect Expression. expected=%q got=%q",
			test, expected[2], node.Right.String())
	}
}

func checkPrefixExpression(
	t *testing.T, test int, expected []string, node *ast.PrefixExpression,
) {
	if node.Operator != expected[0] {
		t.Errorf("tests[%d]: incorrect operator. expected=%q got=%q",
			test, expected[0], node.Operator)
	}

	if node.String() != expected[1] {
		t.Errorf("tests[%d]: incorrect Expression. expected=%q got=%q",
			test, expected[1], node.String())
	}
}

func checkExpressionStatement(
	t *testing.T,
	test int,
	expected []string,
	node *ast.ExpressionStatement,
) {
	switch s := node.Expression.(type) {
	case *ast.Boolean:
		checkBoolean(t, test, expected, s)
	case *ast.Identifier:
		checkIdentifier(t, test, expected, s)
	case *ast.IntegerLiteral:
		checkIntegerLiteral(t, test, expected, s)
	case *ast.PrefixExpression:
		checkPrefixExpression(t, test, expected, s)
	case *ast.InfixExpression:
		checkInfixExpression(t, test, expected, s)
	default:
		t.Errorf("tests[%d]: wrong type. got=%T\n",
			test, node.Expression)
	}
}

func checkDeclarationStatement(
	t *testing.T,
	test int,
	expected []string,
	stmt *ast.DeclarationStatement,
) {
	// int
	if expected[0] != stmt.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect type. expected=%q got=%q",
			test, expected[0], stmt.TokenLiteral())
	}

	// x
	if expected[1] != stmt.Name.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect identifier. expected=%q got=%q",
			test, expected[1], stmt.Name.TokenLiteral())
	}

	// =

	// y
	if expected[2] != stmt.Value.String() {
		t.Errorf("tests[%d]: incorrect identifier. expected=%q got=%q",
			test, expected[1], stmt.Name.TokenLiteral())
	}

	// ;
}

func checkReturnStatement(
	t *testing.T, test int, expected []string, stmt *ast.ReturnStatement,
) {
	// return
	if expected[0] != stmt.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect type. expected=%q got=%q",
			test, expected[0], stmt.TokenLiteral())
	}

	// x
	if expected[1] != stmt.ReturnValue.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect type. expected=%q got=%q",
			test, expected[0], stmt.ReturnValue.TokenLiteral())
	}
}

// ----------------------------------------------------------------------------
// Statement tests
// ----------------------------------------------------------------------------

func TestDeclarationStatement(t *testing.T) {
	input := `
	int x = 0;
	int y = x;
	`
	expected := testResults{{"int", "x", "0"}, {"int", "y", "x"}}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

func TestReturnStatement(t *testing.T) {
	input := "return x;"
	expected := testResults{{"return", "x"}}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

// ----------------------------------------------------------------------------
// Expression tests
// ----------------------------------------------------------------------------

func TestBangOperator(t *testing.T) {
	input := "!true; !false; !!true; !!false;"
	expected := testResults{
		{"!", "(!true)"},
		{"!", "(!false)"},
		{"!", "(!(!true))"},
		{"!", "(!(!false))"},
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

func TestBooleanExpressions(t *testing.T) {
	input := "true; false;"
	expected := testResults{{"true"}, {"false"}}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

func TestAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x = 2;", "(x = 2)"},
		{"foo = 2 * y;", "(foo = (2 * y))"},
		{"foo = bar = baz = 100;", "(foo = (bar = (baz = 100)))"},
		{"foo = bar = baz = x + 3;", "(foo = (bar = (baz = (x + 3))))"},
	}

	for index, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		checkProgram(t, program)
		checkErrors(t, p)
		checkLength(t, 1, program.Statements)

		for _, statement := range program.Statements {
			if statement.String() != test.expected {
				t.Errorf(
					`tests[%d]: parser tree incorrect.
					expected=%q got=%q`,
					index,
					test.expected,
					statement.String())
			}
		}
	}
}

func TestParentheses(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"4 * (2 + 3);", "(4 * (2 + 3))"},
		{"(2 + 3) * 4;", "((2 + 3) * 4)"},
	}

	for index, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		checkProgram(t, program)
		checkErrors(t, p)
		checkLength(t, 1, program.Statements)

		for _, statement := range program.Statements {
			if statement.String() != test.expected {
				t.Errorf(`tests[%d]: parser tree incorrect.
					expected=%q got=%q`,
					index,
					test.expected,
					statement.String())
			}
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	expected := testResults{{"foobar"}}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	expected := testResults{{"5"}}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

func TestEqualityExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     interface{}
		right    interface{}
		operator string
	}{
		{"true == true;", true, true, "=="},
		{"false == false;", false, false, "=="},
		{"5 == true;", int64(5), true, "=="},
		{"foo == bar;", "foo", "bar", "=="},

		{"true != true;", true, true, "!="},
		{"false != false;", false, false, "!="},
		{"5 != true;", int64(5), true, "!="},
		{"foo != bar;", "foo", "bar", "!="},

		{"true < true;", true, true, "<"},
		{"false < false;", false, false, "<"},
		{"5 < true;", int64(5), true, "<"},
		{"foo < bar;", "foo", "bar", "<"},

		{"true <= true;", true, true, "<="},
		{"false <= false;", false, false, "<="},
		{"5 <= true;", int64(5), true, "<="},
		{"foo <= bar;", "foo", "bar", "<="},

		{"true > true;", true, true, ">"},
		{"false > false;", false, false, ">"},
		{"5 > true;", int64(5), true, ">"},
		{"foo > bar;", "foo", "bar", ">"},

		{"true >= true;", true, true, ">="},
		{"false >= false;", false, false, ">="},
		{"5 >= true;", int64(5), true, ">="},
		{"foo >= bar;", "foo", "bar", ">="},
	}

	for index, test := range tests {
		l := lexer.New(test.input)
		p := New(l)
		program := p.ParseProgram()

		checkProgram(t, program)
		checkErrors(t, p)
		checkLength(t, 1, program.Statements)

		for _, statement := range program.Statements {
			es, ok := statement.(*ast.ExpressionStatement)
			if !ok {
				t.Errorf(`tests[%d]: expected
					ast.ExpressionStatement got=%T (%+v)`,
					index, statement, statement)
			}

			ie, ok := es.Expression.(*ast.InfixExpression)
			if !ok {
				t.Errorf(`tests[%d]: expected 
					ast.InfixExpression got=%T (%+v)`,
					index, es, es)
			}

			if ie.Operator != test.operator {
				t.Errorf(`tests[%d]: operator wrong. 
					expected=%q got=%q (%+v)`,
					index, test.operator, ie.Operator, ie)
			}

			switch testLeft := test.left.(type) {
			case bool:
				ieLeft, ok := ie.Left.(*ast.Boolean)
				if !ok {
					t.Errorf(`tests[%d]: left operand type 
						wrong. expected=Boolean got=%T 
					(%+v)`, index, ie.Left, ie)
				} else {
					if ieLeft.Value != testLeft {
						t.Errorf(`tests[%d]: left value
							wrong. expected=%t 
							got=%t (%+v)`,
							index,
							testLeft,
							ieLeft.Value,
							ieLeft)
					}
				}
			case int64:
				ieLeft, ok := ie.Left.(*ast.IntegerLiteral)
				if !ok {
					t.Errorf(`tests[%d]: left operand type 
						wrong. expected=IntegerLiteralT
					got=%T (%+v)`, index, ie.Left, ie)
				} else {
					if ieLeft.Value != testLeft {
						t.Errorf(`tests[%d]: left value
							wrong. expected=%d 
							got=%d (%+v)`,
							index, testLeft,
							ieLeft.Value, ieLeft)
					}
				}
			case string:
				ieLeft, ok := ie.Left.(*ast.Identifier)
				if !ok {
					t.Errorf(`tests[%d]: left operand type
						wrong. expected=Identifier 
					got=%T (%+v)`, index, ie.Left, ie)
				} else {
					if ieLeft.Value != testLeft {
						t.Errorf(`tests[%d]: left value
							wrong. expected=%s 
							got=%s (%+v)`,
							index,
							testLeft,
							ieLeft.Value,
							ieLeft)
					}
				}

			default:
				t.Errorf(`tests[%d]: left type unsupported. 
					left=%T (%+v)`, index, test.left, test)
			}

			switch testRight := test.right.(type) {
			case bool:
				ieRight, ok := ie.Right.(*ast.Boolean)
				if !ok {
					t.Errorf(`tests[%d]: right operand type
						wrong. expected=%T 
						got=%T (%+v)`,
						index,
						test.right,
						ie.Right,
						ie)
				} else {
					if ieRight.Value != testRight {
						t.Errorf(`tests[%d]: right
							value wrong. 
							expected=%t 
							got=%t (%+v)`,
							index,
							testRight,
							ieRight.Value,
							ieRight)
					}
				}
			case int64:
				ieRight, ok := ie.Right.(*ast.IntegerLiteral)
				if !ok {
					t.Errorf(`tests[%d]: right operand type
						wrong. expected=%T 
						got=%T (%+v)`,
						index,
						test.right,
						ie.Right,
						ie)
				} else {
					if ieRight.Value != testRight {
						t.Errorf(`tests[%d]: right 
							value wrong. 
							expected=%d 
							got=%d (%+v)`,
							index,
							testRight,
							ieRight.Value,
							ieRight)
					}
				}
			case string:
				ieRight, ok := ie.Right.(*ast.Identifier)
				if !ok {
					t.Errorf(`tests[%d]: right operand type
						wrong. expected=%T 
						got=%T (%+v)`,
						index,
						test.right,
						ie.Right,
						ie)
				} else {
					if ieRight.Value != testRight {
						t.Errorf(`tests[%d]: right 
							value wrong. 
							expected=%s 
							got=%s (%+v)`,
							index,
							testRight,
							ieRight.Value,
							ieRight)
					}
				}

			default:
				t.Errorf(`tests[%d]: right type unsupported. 
					left=%T (%+v)`,
					index, test.right, test)
			}
		}
	}
}

func TestArithmeticExpressions(t *testing.T) {
	input := `
		 10 +  10;
		-10 +  10;
		 10 + -10;
		-10 + -10;

		 10 -  10;
		-10 -  10;
		 10 - -10;
		-10 - -10;

		 15 *  15;
		-15 *  15;
		 15 * -15;
		-15 * -15;

		 15 /  15;
		-15 /  15;
		 15 / -15;
		-15 / -15;
		`
	expected := testResults{
		{"10", "+", "10"},
		{"(-10)", "+", "10"},
		{"10", "+", "(-10)"},
		{"(-10)", "+", "(-10)"},

		{"10", "-", "10"},
		{"(-10)", "-", "10"},
		{"10", "-", "(-10)"},
		{"(-10)", "-", "(-10)"},

		{"15", "*", "15"},
		{"(-15)", "*", "15"},
		{"15", "*", "(-15)"},
		{"(-15)", "*", "(-15)"},

		{"15", "/", "15"},
		{"(-15)", "/", "15"},
		{"15", "/", "(-15)"},
		{"(-15)", "/", "(-15)"},
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

func TestInfixOperatorPrecedence(t *testing.T) {
	input := `
		  5 +  10 * -15;
		-10 +  15 *   5;
		 15 +  -5 *  10;
		 -5 + -10 * -15;

		  5 -  10 * -15;
		-10 -  15 *   5;
		 15 -  -5 *  10;
		 -5 - -10 * -15;

		  5 *  10 + -15;
		-10 *  15 +   5;
		 15 *  -5 - -10;
		 -5 * -10 -  15;
		`
	expected := testResults{
		{"5", "+", "(10 * (-15))"},
		{"(-10)", "+", "(15 * 5)"},
		{"15", "+", "((-5) * 10)"},
		{"(-5)", "+", "((-10) * (-15))"},

		{"5", "-", "(10 * (-15))"},
		{"(-10)", "-", "(15 * 5)"},
		{"15", "-", "((-5) * 10)"},
		{"(-5)", "-", "((-10) * (-15))"},

		{"(5 * 10)", "+", "(-15)"},
		{"((-10) * 15)", "+", "5"},
		{"(15 * (-5))", "-", "(-10)"},
		{"((-5) * (-10))", "-", "15"},
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

func TestPrefixOperatorExpressions(t *testing.T) {
	input := `
		-10;
		--5;
		`
	expected := testResults{
		{"-", "(-10)"},
		{"-", "(-(-5))"},
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, len(expected), program.Statements)
	checkStatements(t, expected, program.Statements)
}

// ----------------------------------------------------------------------------
// Miscellaneous tests
// ----------------------------------------------------------------------------

func TestProgramString(t *testing.T) {
	input := "int x = 5;"
	expected := "int x = 5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkProgram(t, program)
	checkErrors(t, p)

	if program.String() != expected {
		t.Errorf("program.String() returned %q, expected %q\n",
			program.String(), expected)
	}
}
