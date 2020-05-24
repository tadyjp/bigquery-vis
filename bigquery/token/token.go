package token

type TokenType string

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + literals
	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	NUMBER = "NUMBER" // 1.2
	STRING = "STRING"

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

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	SELECT = "SELECT"
	FROM   = "FROM"
	TRUE   = "TRUE"
	FALSE  = "FALSE"
	NULL   = "NULL"
	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
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
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
