%{
package bigquery

import (
	"strings"
	"text/scanner"

	"github.com/k0kubun/pp"
)

type Token struct {
	tok int
	lit string
	pos scanner.Position
}

%}

%union{
    tok Token
    stmts Statements
    stmt Statement
    expr  Expression
}

%type<expr> program
%type<stmts> statements
%type<stmt> statement
%type<expr> expr
%token<tok> NUMBER

%left '+' '-'
%left '*' '/'

%%

program
    : statements
    {
        $$ = $1
        yylex.(*Lexer).program = $$
    }

statements
    : statement
    {
        $$ = []Statement{$1}
    }
    | statements statement
    {
        $$ = append($1, $2)
    }

statement
    : expr
    {
        $$ = ExpressionStmt{Expr: $1}
    }
    | expr ';'
    {
        $$ = ExpressionStmt{Expr: $1}
    }

expr
    : NUMBER
    {
        $$ = NumExpr{lit: $1.lit}
    }
    | expr '+' expr
    {
        $$ = BinOpExpr{left: $1, operator: '+', right: $3}
    }
    | expr '-' expr
    {
        $$ = BinOpExpr{left: $1, operator: '-', right: $3}
    }
    | expr '*' expr
    {
        $$ = BinOpExpr{left: $1, operator: '*', right: $3}
    }
    | expr '/' expr
    {
        $$ = BinOpExpr{left: $1, operator: '/', right: $3}
    }

%%

type Lexer struct {
	scanner.Scanner
	recentLit  string
	recentPos  scanner.Position
	program Expression
}

func (l *Lexer) Lex(lval *yySymType) int {
	tok := int(l.Scanner.Scan())
	if tok == scanner.Int {
        tok = NUMBER
    }
	lit := l.Scanner.TokenText()
	pos := l.Scanner.Pos()
	lval.tok = Token{tok: tok, lit: lit, pos: pos}
	pp.Println(lval)
	l.recentLit = lit
	l.recentPos = pos
	return tok
}

func (l *Lexer) Error(e string) {
    pp.Println(l)
    panic(e)
}

func Parse(s string) Expression {
	l := new(Lexer)
    l.Scanner.Init(strings.NewReader(s))
	yyParse(l)
	return l.program
}
