package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/freddiehaddad/corrosion/pkg/lexer"
	"github.com/freddiehaddad/corrosion/pkg/token"
)

const appName = "Corrosion"
const prompt = "> "

func main() {
	fmt.Println("Welcome to", appName)
	fmt.Println("")
	fmt.Println("Press Ctrl+D (^D) to exit")

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for scanner.Scan() {
		input := scanner.Text()

		l := lexer.New(input)
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}

		fmt.Print(prompt)
	}

	fmt.Println("Exiting", appName)
}
