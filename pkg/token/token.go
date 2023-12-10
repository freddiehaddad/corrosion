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

	// operators
	ASSIGN   = "="
	MINUS    = "-"
	PLUS     = "+"
	DIVIDE   = "/"
	MULTIPLY = "*"

	// keywords
	INT    = "INT"
	RETURN = "RETURN"

	TRUE  = "TRUE"
	FALSE = "FALSE"

	// literals
	IDENT   = "IDENT"
	INTEGER = "INTEGER"

	// end of input
	EOF       = "EOF"
	EOF_VALUE = byte(0)

	// unsupported input
	ILLEGAL = "ILLEGAL"
)

// Language keywords.
var keywords = map[string]TokenType{
	"int":    INT,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

// Checks if tt is in the keyword table and return the corresponding TokenType.
// Otherwise, IDENT is returned.
func LookupType(tt string) TokenType {
	if tok, ok := keywords[tt]; ok {
		return tok
	}
	return IDENT
}
