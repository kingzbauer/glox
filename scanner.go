package main

import (
	"strconv"
)

// Scanner type
type Scanner struct {
	source               string
	tokens               []Token
	start, current, line int
	lox                  *Lox
}

// NewScanner instance from source
func NewScanner(source string, lox *Lox) *Scanner {
	return &Scanner{
		source: source,
		tokens: []Token{},
		line:   1,
		lox:    lox,
	}
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	lexeme := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tokenType, lexeme, literal, s.line))
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LeftParen, nil)
	case ')':
		s.addToken(RightParen, nil)
	case '{':
		s.addToken(LeftBrace, nil)
	case '}':
		s.addToken(RightBrace, nil)
	case ',':
		s.addToken(Comma, nil)
	case '.':
		s.addToken(Dot, nil)
	case '-':
		s.addToken(Minus, nil)
	case '+':
		s.addToken(Plus, nil)
	case ';':
		s.addToken(SemiColon, nil)
	case '*':
		s.addToken(Star, nil)
	case '!':
		if s.match('=') {
			s.addToken(BangEqual, nil)
		} else {
			s.addToken(Bang, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EqualEqual, nil)
		} else {
			s.addToken(Equal, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GreaterEqual, nil)
		} else {
			s.addToken(Greater, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LessEqual, nil)
		} else {
			s.addToken(Less, nil)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.blockComments()
		} else {
			s.addToken(Slash, nil)
		}
	case ' ':
		fallthrough
	case '\t':
		fallthrough
	case '\r':
	case '\n':
		s.line++
	case '"':
		s.str()
	default:
		if s.isDigit(c) {
			s.number()
		} else if s.isAlpha(c) {
			s.identifier()
		} else {
			// handle error
			s.lox.Exception(s.line, "Unexpected character")
		}
	}
}

// ScanTokens output
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) str() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.lox.Exception(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the sorrounding quotes.
	value := s.source[s.start+1 : s.current-1]
	s.addToken(String, value)
}

func (s *Scanner) isDigit(val byte) bool {
	return val >= '0' && val <= '9'
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	// We will ignore the error. We expect the input to be valid
	number, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
	s.addToken(Number, number)
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') ||
		(char >= 'A' && char <= 'Z') ||
		char == '_'
}

func (s *Scanner) isAlphaNumeric(char byte) bool {
	return s.isAlpha(char) || s.isDigit(char)
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	typ, exists := Keywords[text]
	if !exists {
		typ = Identifier
	}
	s.addToken(typ, nil)
}

func (s *Scanner) blockComments() {
	for !s.isAtEnd() {
		char := s.advance()
		if char == '*' && s.match('/') {
			return
		}
		if char == '\n' {
			s.line++
		}
	}
	s.lox.Exception(s.line, "Unterminated block comment")
}
