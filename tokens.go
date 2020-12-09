package main

import (
	"fmt"
)

// TokenType defines the type of lexeme
type TokenType int

const (
	// Single character tokens

	// LeftParen token type
	LeftParen TokenType = iota
	// RightParen token type
	RightParen
	// LeftBrace token type
	LeftBrace
	// RightBrace token type
	RightBrace
	// Comma token type
	Comma
	// Dot token type
	Dot
	// Minus token type
	Minus
	// Plus token type
	Plus
	// SemiColon token type
	SemiColon
	// Slash token type
	Slash
	// Start token type
	Start

	// One or two character tokens

	// Bang token type
	Bang
	// BangEqual token type
	BangEqual
	// Equal token type
	Equal
	// EqualEqual token type
	EqualEqual
	// Greater token type
	Greater
	// GreaterEqual token type
	GreaterEqual
	// Less token type
	Less
	// LessEqual token type
	LessEqual

	// Literals
	Identifier
	String
	Number

	// Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

// Token instance
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// NewToken instance
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

// String method
func (t Token) String() string {
	return fmt.Sprintf("%d %s %s", t.Type, t.Lexeme, t.Literal)
}
