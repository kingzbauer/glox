package main

/*
	expression     → literal
								 | unary
								 | binary
								 | grouping ;

	literal        → NUMBER | STRING | "true" | "false" | "nil" ;
	grouping       → "(" expression ")" ;
	unary          → ( "-" | "!" ) expression ;
	binary         → expression operator expression ;
	operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
								 | "+"  | "-"  | "*" | "/" ;
*/

// Expr type
type Expr interface {
	Accept(Visitor) interface{}
}

// Binary type of expression
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

// Accept Expr interface
func (b Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinary(b)
}

// Grouping type of expression
type Grouping struct {
	Expr Expr
}

// Accept Expr interface
func (g Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGrouping(g)
}

// Literal type of expression
type Literal struct {
	Value interface{}
}

// Accept Expr interface
func (l Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(l)
}

// Unary type of expression
type Unary struct {
	Operator Token
	Right    Expr
}

// Accept Expr interface
func (u Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnary(u)
}
