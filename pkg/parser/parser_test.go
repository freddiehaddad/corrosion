package parser

import (
	"testing"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/lexer"
)

type decTest struct {
	decType string
	decName string
}

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

func checkLength(t *testing.T, tests []decTest, stmts []ast.Statement) {
	if len(stmts) != len(tests) {
		t.Fatalf("ParseProgram returned %d statements, expected %d\n", len(stmts), len(tests))
	}
}

func checkStatements(t *testing.T, tests []decTest, stmts []ast.Statement) {
	for i, test := range tests {
		stmt, ok := stmts[i].(*ast.DeclarationStatement)
		if !ok {
			t.Errorf("tests[%d]: wrong type. expected=DeclarationStatement got=%T\n", i, stmt)
			continue
		}

		checkStatement(t, i, test, stmt)
	}
}

func checkStatement(t *testing.T, test int, expected decTest, stmt *ast.DeclarationStatement) {
	if expected.decType != stmt.TokenLiteral() {
		t.Errorf("tests[%d]: incorrect type. expected=%q got=%q\n",
			test, expected.decType, stmt.TokenLiteral())
	}

	if expected.decName != stmt.Name.Value {
		t.Errorf("tests[%d]: incorrect identifier. expected=%q got=%q\n",
			test, expected.decName, stmt.Name.Value)
	}
}

func TestDeclarationStatement(t *testing.T) {
	input := `
		int x = 0;
		int y = 10;
		int bigNumber = 999999;
		`

	tests := []decTest{
		{"int", "x"},
		{"int", "y"},
		{"int", "bigNumber"},
	}

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	checkProgram(t, program)
	checkErrors(t, p)
	checkLength(t, tests, program.Statements)
	checkStatements(t, tests, program.Statements)
}
