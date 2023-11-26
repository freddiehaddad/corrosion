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
		default:
			t.Errorf("tests[%d]: wrong type. got=%T\n", i, stmt)
		}
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
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
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
