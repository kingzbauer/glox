package main

import (
	"errors"
	"fmt"
)

/*
	expression     → equality ;
	equality       → comparison ( ( "!=" | "==" ) comparison )* ;
	comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
	term           → factor ( ( "-" | "+" ) factor )* ;
	factor         → unary ( ( "/" | "*" ) unary )* ;
	unary          → ( "!" | "-" ) unary
								 | primary ;
	primary        → NUMBER | STRING | "true" | "false" | "nil"
								 | "(" expression ")" ;
*/

// ErrParse error
var ErrParse = errors.New("Syntax error")

// Parser type
type Parser struct {
	tokens  []Token
	current int
	lox     *Lox
}

// NewParser factory
func NewParser(tokens []Token, lox *Lox) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		lox:     lox,
	}
}

// Expression func
func (p *Parser) Expression() Expr {
	return p.Equality()
}

// Equality func
func (p *Parser) Equality() Expr {
	expr := p.Comparison()

	for p.match(BangEqual, EqualEqual) {
		operator := p.previous()
		right := p.Comparison()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// Comparison func
func (p *Parser) Comparison() Expr {
	expr := p.Term()

	for p.match(Greater, GreaterEqual, Less, LessEqual) {
		operator := p.previous()
		right := p.Term()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// Term func
func (p *Parser) Term() Expr {
	expr := p.Factor()

	for p.match(Minus, Plus) {
		operator := p.previous()
		right := p.Factor()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// Factor func
func (p *Parser) Factor() Expr {
	expr := p.Unary()

	for p.match(Slash, Star) {
		operator := p.previous()
		right := p.Unary()
		expr = Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

// Unary func
func (p *Parser) Unary() Expr {
	if p.match(Bang, Minus) {
		operator := p.previous()
		right := p.Unary()
		return Unary{Operator: operator, Right: right}
	}

	return p.Primary()
}

// Primary func
func (p *Parser) Primary() Expr {
	token := p.peek()
	switch t := token.Type; t {
	case False:
		p.advance()
		return Literal{Value: false}
	case True:
		p.advance()
		return Literal{Value: true}
	case Nil:
		p.advance()
		return Literal{Value: nil}
	case Number, String:
		return Literal{Value: p.advance().Literal}
	case LeftParen:
		p.advance()
		expr := p.Expression()
		p.consume(RightParen, "Expect ')' after expression.")
		return Grouping{Expr: expr}
	default:
		panic(p.err(p.peek(), "Expected expression."))
	}
}

func (p *Parser) consume(tokenType TokenType, errMsg string) Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(p.err(p.peek(), errMsg))
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == tokenType
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) err(token Token, msg string) error {
	if token.Type == EOF {
		p.lox.report(token.Line, "at end", msg)
	} else {
		p.lox.report(token.Line, fmt.Sprintf("at '%s'", token.Lexeme), msg)
	}

	return ErrParse
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SemiColon {
			return
		}

		switch t := p.peek().Type; t {
		case Class:
			fallthrough
		case Fun:
			fallthrough
		case Var:
			fallthrough
		case For:
			fallthrough
		case If:
			fallthrough
		case While:
			fallthrough
		case Print:
			fallthrough
		case Return:
			return
		}
		p.advance()
	}
}

// Parse starts the parsing process
func (p *Parser) Parse() (expr Expr) {
	defer func() {
		if err := recover(); err != nil {
			if recErr, ok := err.(error); ok && recErr == ErrParse {
				expr = nil
			} else {
				panic(err)
			}
		}
	}()

	expr = p.Expression()
	return
}
