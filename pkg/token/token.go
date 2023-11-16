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
	SEMICOLON = "SEMICOLON"
	COMMA     = "COMMA"
	LPAREN    = "LEFT PAREN"
	RPAREN    = "RIGHT PAREN"
	LBRACKET  = "LEFT BRACKET"
	RBRACKET  = "RIGHT BRACKET"
	LBRACE    = "LEFT BRACE"
	RBRACE    = "RIGHT BRACE"

	// operators
	BITWISE_AND  = "BITWISE AND"
	BITWISE_OR   = "BITWISE OR"
	BITWISE_XOR  = "BITWISE XOR"
	BITWISE_NOT  = "BITWISE NOT"
	BANG         = "BANG"
	ASSIGN       = "ASSIGN"
	PLUS         = "PLUS"
	MINUS        = "MINUS"
	ASTERISK     = "ASTERISK"
	FSLASH       = "FORWARD SLASH"
	LESS_THAN    = "LESS THAN"
	GREATER_THAN = "GREATER THAN"

	// logical operators
	EQUAL     = "EQUAL"
	NOT_EQUAL = "NOT EQUAL"
	AND       = "LOGICAL AND"
	OR        = "LOGICAL OR"

	// end of input
	EOF       = "EOF"
	EOF_VALUE = byte(0)

	// unknown input
	UNKNOWN = "UNKNOWN"
)
