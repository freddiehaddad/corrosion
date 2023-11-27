package parser

import (
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
		t.Fatalf("ParseProgram returned %d statements, expected %d\n", len(stmts), length)
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

func checkIdentifier(t *testing.T, test int, expected []string, node *ast.Identifier) {
	if node.Value != expected[0] {
		t.Errorf("tests[%d]: incorrect value. expected=%q got=%q\n",
			test, expected[0], node.Value)
	}
}

func checkIntegerLiteral(t *testing.T, test int, expected []string, node *ast.IntegerLiteral) {
	if node.Value != expected[0] {
		t.Errorf("tests[%d]: incorrect value. expected=%q got=%q\n",
			test, expected[0], node.Value)
	}
}

func checkInfixExpression(t *testing.T, test int, expected []string, node *ast.InfixExpression) {
	if node.Left.String() != expected[0] {
		t.Errorf("tests[%d]: incorrect Expression. expected=%q got=%q\n",
			test, expected[0], node.Left.String())
	}

	if node.Operator != expected[1] {
		t.Errorf("tests[%d]: incorrect operator. expected=%q got=%q\n",
			test, expected[1], node.Operator)
	}

	if node.Right.String() != expected[2] {
		t.Errorf("tests[%d]: incorrect Expression. expected=%q got=%q\n",
			test, expected[2], node.Right.String())
	}
}

func checkPrefixExpression(t *testing.T, test int, expected []string, node *ast.PrefixExpression) {
	if node.Operator != expected[0] {
		t.Errorf("tests[%d]: incorrect operator. expected=%q got=%q\n",
			test, expected[0], node.Operator)
	}

	if node.Right.String() != expected[1] {
		t.Errorf("tests[%d]: incorrect Expression. expected=%q got=%q\n",
			test, expected[1], node.Right.String())
	}
}

func checkExpressionStatement(t *testing.T, test int, expected []string, node *ast.ExpressionStatement) {
	switch s := node.Expression.(type) {
	case *ast.Identifier:
		checkIdentifier(t, test, expected, s)
	case *ast.IntegerLiteral:
		checkIntegerLiteral(t, test, expected, s)
	case *ast.PrefixExpression:
		checkPrefixExpression(t, test, expected, s)
	case *ast.InfixExpression:
		checkInfixExpression(t, test, expected, s)
	default:
		t.Errorf("tests[%d]: wrong type. got=%T\n", test, node.Expression)
	}
}

func checkDeclarationStatement(t *testing.T, test int, expected []string, stmt *ast.DeclarationStatement) {
	// int
	if expected[0] != stmt.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect type. expected=%q got=%q\n",
			test, expected[0], stmt.TokenLiteral())
	}

	// x
	if expected[1] != stmt.Name.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect identifier. expected=%q got=%q\n",
			test, expected[1], stmt.Name.TokenLiteral())
	}

	// =

	// y
	if expected[2] != stmt.Value.String() {
		t.Errorf("tests[%d]: incorrect identifier. expected=%q got=%q\n",
			test, expected[1], stmt.Name.TokenLiteral())
	}

	// ;
}

func checkReturnStatement(t *testing.T, test int, expected []string, stmt *ast.ReturnStatement) {
	// return
	if expected[0] != stmt.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect type. expected=%q got=%q\n",
			test, expected[0], stmt.TokenLiteral())
	}

	// x
	if expected[1] != stmt.ReturnValue.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect type. expected=%q got=%q\n",
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

func TestInfixOperatorExpressions(t *testing.T) {
	{
		input := `
			 10 +  10;
			-10 +  10;
			 10 + -10;
			-10 + -10;
			 10 -  10;
			-10 -  10;
			 10 - -10;
			-10 - -10;
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
		}

		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()

		checkProgram(t, program)
		checkErrors(t, p)
		checkLength(t, len(expected), program.Statements)
		checkStatements(t, expected, program.Statements)
	}
}

func TestPrefixOperatorExpressions(t *testing.T) {
	{
		input := `
			-10;
			!5;
		`
		expected := testResults{
			{"-", "10"},
			{"!", "5"},
		}

		l := lexer.New(input)
		p := New(l)
		program := p.ParseProgram()

		checkProgram(t, program)
		checkErrors(t, p)
		checkLength(t, len(expected), program.Statements)
		checkStatements(t, expected, program.Statements)
	}
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
		t.Errorf("program.String() returned %q, expected %q\n", program.String(), expected)
	}
}
