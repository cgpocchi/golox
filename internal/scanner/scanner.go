// Package for converting a source program into a sequence of tokens.
package scanner

import (
	"golox/internal/lox"
	"golox/internal/token"
	"golox/internal/utils"
	"strconv"
)

// Struct for scanning a source string into a slice of tokens.
type Scanner struct {
	source     string
	tokens     []token.Token
	errTracker *lox.ErrorTracker
	start      int
	current    int
	line       int
}

var keywords = map[string]token.TokenType{
	"and":    token.AND,
	"class":  token.CLASS,
	"else":   token.ELSE,
	"false":  token.FALSE,
	"for":    token.FOR,
	"fun":    token.FUN,
	"if":     token.IF,
	"none":   token.NONE,
	"or":     token.OR,
	"print":  token.PRINT,
	"return": token.RETURN,
	"super":  token.SUPER,
	"this":   token.THIS,
	"true":   token.TRUE,
	"var":    token.VAR,
	"while":  token.WHILE,
}

func NewScanner(source string, errTracker *lox.ErrorTracker) *Scanner {
	return &Scanner{
		source:     source,
		tokens:     make([]token.Token, 0, len(source)),
		errTracker: errTracker,
		start:      0,
		current:    0,
		line:       1,
	}
}

// Convert the source string into a slice of tokens.
func (s *Scanner) ScanTokens() []token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.Token{
		Type:    token.EOF,
		Lexeme:  "",
		Literal: nil,
		Line:    s.line,
	})
	return s.tokens
}

// Scan for the next token in the source string and add it to the token list.
func (s *Scanner) scanToken() {
	switch c := s.advance(); c {
	case '(':
		s.addTokenNoLit(token.LEFT_PAREN)
	case ')':
		s.addTokenNoLit(token.RIGHT_PAREN)
	case '{':
		s.addTokenNoLit(token.LEFT_BRACE)
	case '}':
		s.addTokenNoLit(token.RIGHT_BRACE)
	case ',':
		s.addTokenNoLit(token.COMMA)
	case '.':
		s.addTokenNoLit(token.DOT)
	case '-':
		s.addTokenNoLit(token.MINUS)
	case '+':
		s.addTokenNoLit(token.PLUS)
	case ';':
		s.addTokenNoLit(token.SEMICOLON)
	case '*':
		s.addTokenNoLit(token.STAR)
	case '!':
		s.addTokenNoLit(utils.TernaryOp(s.nextCharMatches('='), token.BANG_EQUAL, token.BANG))
	case '=':
		s.addTokenNoLit(utils.TernaryOp(s.nextCharMatches('='), token.EQUAL_EQUAL, token.EQUAL))
	case '<':
		s.addTokenNoLit(utils.TernaryOp(s.nextCharMatches('='), token.LESS_EQUAL, token.LESS))
	case '>':
		s.addTokenNoLit(utils.TernaryOp(s.nextCharMatches('='), token.GREATER_EQUAL, token.GREATER))
	case '/':
		// ignore rest of the line since it is a comment
		if s.nextCharMatches('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addTokenNoLit(token.SLASH)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		s.addStringToken()
	default:
		if isDigit(c) {
			s.addNumberToken()
		} else if isAlpha(c) {
			s.addIdentifierToken()
		} else {
			s.errTracker.Error(s.line, "Unexpected character.")
		}
	}
}

// Helper to determine if the given byte is a letter or digit.
func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}

// Helper to check if the given byte is an alphabetical character or wild card.
func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

// Helper to check if the given byte is a digit.
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// Get but don't consume the next character.
func (s *Scanner) peekNext() byte {
	if s.current+1 > len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

// Add identifier (variable name or keyword) token.
func (s *Scanner) addIdentifierToken() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]

	if tokenType, ok := keywords[text]; !ok {
		s.addTokenNoLit(token.IDENTIFIER)
	} else {
		s.addTokenNoLit(tokenType)
	}
}

// Add number literal token.
func (s *Scanner) addNumberToken() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// look for fractional part
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// consume the decimal
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}
	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.errTracker.Error(s.line, "Error parsing number")
	} else {
		s.addToken(token.NUMBER, num)
	}
}

// Add string literal token.
func (s *Scanner) addStringToken() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.errTracker.Error(s.line, "Unterminated string.")
		return
	}

	// consume end quote
	s.advance()

	// Trim the surrounding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

// Get the next character in the source string without consuming it.
func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

// Check that next character in sources matches the given character.
// If next character matches consume it and return true. Otherwise return false.
func (s *Scanner) nextCharMatches(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

// Consume the next character in source string.
func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

// Add a token of given type without a literal value.
func (s *Scanner) addTokenNoLit(tokenType token.TokenType) {
	s.addToken(tokenType, nil)
}

// Add a token with the given type and literal value.
func (s *Scanner) addToken(tokenType token.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.Token{
		Type:    tokenType,
		Lexeme:  text,
		Literal: literal,
		Line:    s.line,
	})
}

// True if scanner is at the end of the source string, false otherwise.
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}
