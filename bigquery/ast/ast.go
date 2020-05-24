package ast

type Statements []Statement
type Statement interface{}
type Expression interface{}

type ExpressionStmt struct {
	Expr Expression
}

type NumExpr struct {
	lit string
}
type BinOpExpr struct {
	left     Expression
	operator rune
	right    Expression
}
