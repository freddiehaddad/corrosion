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
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	COMMA     = ","

	// operators
	ASSIGN   = "="
	MINUS    = "-"
	PLUS     = "+"
	DIVIDE   = "/"
	MULTIPLY = "*"

	// prefix only operators
	BANG = "!"

	// logical operators
	EQ       = "=="
	NOT_EQ   = "!="
	LT       = "<"
	LT_EQUAL = "<="
	GT       = ">"
	GT_EQUAL = ">="

	// keywords
	VAR    = "VAR"
	FUNC   = "FUNC"
	RETURN = "RETURN"
	IF     = "IF"
	ELSE   = "ELSE"

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
	"var":    VAR,
	"func":   FUNC,
	"if":     IF,
	"else":   ELSE,
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
