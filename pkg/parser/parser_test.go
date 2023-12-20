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
		case *ast.VariableDeclarationStatement:
			checkVariableDeclarationStatement(t, i, expects[i], s)
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

func checkVariableDeclarationStatement(
	t *testing.T,
	test int,
	expected []string,
	stmt *ast.VariableDeclarationStatement,
) {
	// var
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
			test, expected[1], stmt.Value.String())
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

func testInfixExpression(
	t *testing.T, exp ast.Expression, left, op, right string,
) bool {
	ie, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not an ast.InfixExpression got=%T (%+v)",
			exp, exp)
		return false
	}

	if ie.Left.String() != left {
		t.Errorf("left is wrong. expected=%s got=%s", left,
			ie.Left.String())
		return false
	}

	if ie.Operator != op {
		t.Errorf("op is wrong. expected=%s got=%s", op, ie.Operator)
		return false
	}

	if ie.Right.String() != right {
		t.Errorf("right is wrong. expected=%s got=%s",
			right, ie.Right.String())
		return false
	}

	return true
}

func testAssignmentExpression(
	t *testing.T, exp ast.Expression, left, op, right string,
) bool {
	as, ok := exp.(*ast.AssignmentExpression)
	if !ok {
		t.Errorf("expected ast.AssignmentExpression got =%T (%+v)",
			exp, exp)
		return false
	}

	if as.Left.String() != left {
		t.Errorf("left is wrong. expected=%s got=%s",
			left, as.Left.String())
		return false
	}

	if as.Operator != op {
		t.Errorf("op is wrong. expected=%s got=%s", op, as.Operator)
		return false
	}

	if as.Right.String() != right {
		t.Errorf("right is wrong. expected=%s got=%s",
			right, as.Right.String())
		return false
	}

	return true
}

// ----------------------------------------------------------------------------
// Statement tests
// ----------------------------------------------------------------------------

func TestVariableDeclarationStatement(t *testing.T) {
	input := `
	var x = 0;
	var y = x;
	`
	expected := testResults{{"var", "x", "0"}, {"var", "y", "x"}}

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

func TestFunctionDeclaration(t *testing.T) {
	input := `
		func myfunction(foo, bar) {
			var buz = foo + bar;
			return buz;
		}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, 1, program.Statements)

	fs, ok := program.Statements[0].(*ast.FunctionDeclarationStatement)
	if !ok {
		t.Errorf("Expected FunctionDeclarationStatement. got=%T",
			program.Statements[0])
	}

	if fs.Token.Type != "FUNC" {
		t.Errorf("TokenType wrong. expected=%q got=%q",
			"func", fs.Token.Type)
	}

	if fs.Name.Value != "myfunction" {
		t.Errorf("Function name wrong. expected=%q got=%q",
			"myfunction", fs.Name.Value)
	}

	if len(fs.Parameters) != 2 {
		t.Errorf("Wrong number of parameters. Expected=2 Got=%d",
			len(fs.Parameters))
	}

	identifiers := []string{"foo"}
	checkIdentifier(t, 0, identifiers, &fs.Parameters[0])

	identifiers = []string{"bar"}
	checkIdentifier(t, 0, identifiers, &fs.Parameters[1])

	checkLength(t, 2, fs.Body.Statements)

	body := fs.Body.Statements
	test := []string{"var", "buz", "(foo + bar)"}

	ds, ok := body[0].(*ast.VariableDeclarationStatement)
	if !ok {
		t.Errorf("Expected ast.VariableDeclarationStatement. Got=%T",
			body[0])
	}
	checkVariableDeclarationStatement(t, 0, test, ds)

	rs, ok := body[1].(*ast.ReturnStatement)
	if !ok {
		t.Errorf("Expected ast.ReturnStatement. Got=%T", body[1])
	}
	test = []string{"return", "buz"}
	checkReturnStatement(t, 0, test, rs)
}

func TestFunctionCall(t *testing.T) {
	input := "foo(a, 1+1, bar(), foo()());"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, 1, program.Statements)

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("Expected ast.ExpressionStatement. Got=%T",
			program.Statements[0])
	}

	ce, ok := stmt.Expression.(*ast.FunctionCallExpression)
	if !ok {
		t.Errorf("Expected ast.FunctionCallExpression. Got=%T",
			stmt.Expression)
	}

	id, ok := ce.Function.(*ast.Identifier)
	if !ok {
		t.Errorf("Expected ast.Identifier. Got=%T", ce.Function)
	}
	checkIdentifier(t, 0, []string{"foo"}, id)

	// Arguments
	if len(ce.Arguments) != 4 {
		t.Errorf("Arguments length wrong. Expected 4, got %d",
			len(ce.Arguments))
	}

	arg1, ok := ce.Arguments[0].(*ast.Identifier)
	if !ok {
		t.Errorf("Arg1 Expected ast.Identifier, got %T",
			ce.Arguments[0])
	}
	checkIdentifier(t, 0, []string{"a"}, arg1)

	arg2, ok := ce.Arguments[1].(*ast.InfixExpression)
	if !ok {
		t.Errorf("Arg2 Expected ast.InfixExpression, got %T",
			ce.Arguments[1])
	}
	checkInfixExpression(t, 1, []string{"1", "+", "1"}, arg2)

	arg3, ok := ce.Arguments[2].(*ast.FunctionCallExpression)
	if !ok {
		t.Errorf("Arg3 Expected ast.FunctionCallExpression, got %T",
			ce.Arguments[1])
	}
	arg3id, ok := arg3.Function.(*ast.Identifier)
	if !ok {
		t.Errorf("Arg3 Expected ast.Identifer. Got %t", arg3.Function)
	}
	checkIdentifier(t, 0, []string{"bar"}, arg3id)

	if len(arg3.Arguments) != 0 {
		t.Errorf("Arg3 arguments length wrong. Expected 0, got %d",
			len(arg3.Arguments))
	}

	arg4, ok := ce.Arguments[3].(*ast.FunctionCallExpression)
	if !ok {
		t.Errorf("Arg4 Expected ast.FunctionCallExpression, got %T",
			ce.Arguments[3])
	}

	arg4fce, ok := arg4.Function.(*ast.FunctionCallExpression)
	if !ok {
		t.Errorf(`Arg4.Function wrong.
			Expected ast.FunctionCallExpression. Got=%T`,
			arg4.Function)
	}

	arg4fceid, ok := arg4fce.Function.(*ast.Identifier)
	if !ok {
		t.Errorf("Arg4.Function wrong. Expected ast.Identifier. Got=%T",
			arg4fce.Function)
	}
	checkIdentifier(t, 0, []string{"foo"}, arg4fceid)

	if len(arg4fce.Arguments) != 0 {
		t.Errorf("Arg4 arguments length wrong. Expected 0, got %d",
			len(arg4fce.Arguments))
	}

	if len(arg4.Arguments) != 0 {
		t.Errorf("Arg4 arguments length wrong. Expected 0, got %d",
			len(arg4fce.Arguments))
	}
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

func TestParseIfStatement(t *testing.T) {
	input := `
		if (x == 3) {
			y = 2;
		} else {
			y = 4;
		}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, 1, program.Statements)

	is, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Errorf("expected ast.IfStatement got=%T (%+v)",
			program.Statements[0], program.Statements[0])
	}

	if !testInfixExpression(t, is.Condition, "x", "==", "3") {
		t.FailNow()
	}

	checkLength(t, 1, is.Consequence.Statements)
	checkLength(t, 1, is.Alternative.Statements)

	stmt, ok := is.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("expected ast.Consequence.Expressiongot=%T (%+v)",
			is.Consequence.Statements[0],
			is.Consequence.Statements[0])
	}

	if !testAssignmentExpression(t, stmt.Expression, "y", "=", "2") {
		t.FailNow()
	}

	stmt, ok = is.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf(
			"expected ast.Alternative.ExpressionStatement got=%T",
			is.Alternative.Statements[0])
	}

	if !testAssignmentExpression(t, stmt.Expression, "y", "=", "4") {
		t.FailNow()
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

		foo() + 2;
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

		{"foo()", "+", "2"},
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
	input := "var x = 5;"
	expected := "var x = 5;"

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
