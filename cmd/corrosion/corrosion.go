package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/lexer"
	"github.com/freddiehaddad/corrosion/pkg/parser"
)

const appName = "Corrosion"
const prompt = "> "

func evaluate(p *ast.Program) {
	for index, statement := range p.Statements {
		fmt.Printf("Statements[%d]: %s\n", index, statement.TokenLiteral())
	}
}

func checkProgram(p *ast.Program) bool {
	if p != nil {
		return true
	}

	fmt.Printf("ParseProgram returned nil\n")
	return false
}

func checkErrors(p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	fmt.Printf("ParseProgram returned %d errors\n", len(errors))
	for index, error := range errors {
		fmt.Printf("errors[%d]: %s\n", index, error)
	}
}

func main() {
	fmt.Println("Welcome to", appName)
	fmt.Println("")
	fmt.Println("Press Ctrl+D (^D) to exit")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for scanner.Scan() {
		input := scanner.Text()

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		checkErrors(p)

		if checkProgram(program) {
			evaluate(program)
		}

		fmt.Print(prompt)
	}

	fmt.Println("Exiting", appName)
}
