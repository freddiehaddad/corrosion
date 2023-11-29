package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/evaluator"
	"github.com/freddiehaddad/corrosion/pkg/lexer"
	"github.com/freddiehaddad/corrosion/pkg/parser"
)

const appName = "Corrosion"
const prompt = "> "

func evaluate(p *ast.Program) {
	for _, statement := range p.Statements {
		obj := evaluator.Eval(statement)
		fmt.Println(obj.Inspect())
	}
}

// Returns true if there were any parser errors.
func checkAndPrintErrors(p *parser.Parser) bool {
	errors := p.Errors()
	if len(errors) == 0 {
		return false
	}

	fmt.Printf("ParseProgram returned %d errors\n", len(errors))
	for index, error := range errors {
		fmt.Printf("errors[%d]: %s\n", index, error)
	}

	return true
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to", appName)
	fmt.Println("")
	fmt.Println("Press Ctrl+D (^D) to exit")

	fmt.Print(prompt)
	for scanner.Scan() {
		input := scanner.Text()

		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if !checkAndPrintErrors(p) {
			evaluate(program)
		}

		fmt.Print(prompt)
	}

	fmt.Println("Exiting", appName)
}
