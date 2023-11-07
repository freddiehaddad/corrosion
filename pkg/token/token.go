// The token package holds the definitions for all the tokens in the Corrosion
// language.
package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// Tokens
const (
	// delimiters
	SEMICOLON = ";"
	COMMA = ","
	LPAREN    = "("
	RPAREN    = ")"
	LBRACKET  = "["
	RBRACKET  = "]"
	LBRACE    = "{"
	RBRACE    = "}"

	// operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	FSLASH   = "/"

	// end of input
	EOF = "EOF"

	// unknown input
	UNKNOWN = "UNKNOWN"
)
