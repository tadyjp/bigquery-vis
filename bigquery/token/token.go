package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT          = "IDENT"  // add, foobar, x, y, ...
	UINT           = "UINT"   // 123
	NUMBER         = "NUMBER" // 1.2
	STRING_LITERAL = "STRING LITERAL"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	PERIOD    = "."

	// Brackets
	LPAREN  = "("
	RPAREN  = ")"
	LBRACE  = "{"
	RBRACE  = "}"
	LSQUARE = "["
	RSQUARE = "]"
	LANGLE  = "<"
	RANGLE  = ">"

	// Keywords
	SELECT = "SELECT"
	FROM   = "FROM"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	NULL   = "NULL"
	ARRAY  = "ARRAY"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"

	// Types
	STRING = "STRING"
)

type Token struct {
	Type    TokenType
	Start   int
	Literal string
	Line    int
	Column  int
}

var keywords = map[string]TokenType{
	"select": SELECT,
	"from":   FROM,
	"true":   TRUE,
	"false":  FALSE,
	"null":   NULL,
	"array":  ARRAY,
	"string": STRING,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
