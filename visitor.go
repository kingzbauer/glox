package main

import (
	"fmt"
	"strings"
)

// Visitor type
type Visitor interface {
	VisitBinary(Binary) interface{}
	VisitGrouping(Grouping) interface{}
	VisitLiteral(Literal) interface{}
	VisitUnary(Unary) interface{}
}

// AstPrinter not so-pretty printer
type AstPrinter interface {
	Print(Expr) string
}

// NewAstPrinter factory
func NewAstPrinter() AstPrinter {
	return astPrinter{}
}

type astPrinter struct{}

func (ast astPrinter) Print(expr Expr) string {
	return expr.Accept(ast).(string)
}

func (ast astPrinter) parenthesize(name string, expr ...Expr) string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(name)

	for _, exp := range expr {
		b.WriteString(" ")
		b.WriteString(exp.Accept(ast).(string))
	}
	b.WriteString(")")

	return b.String()
}

func (ast astPrinter) VisitBinary(b Binary) interface{} {
	return ast.parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (ast astPrinter) VisitLiteral(l Literal) interface{} {
	if l.Value == nil {
		return "nil"
	}

	return fmt.Sprint(l.Value)
}

func (ast astPrinter) VisitGrouping(g Grouping) interface{} {
	return ast.parenthesize("group", g.Expr)
}

func (ast astPrinter) VisitUnary(u Unary) interface{} {
	return ast.parenthesize(u.Operator.Lexeme, u.Right)
}
