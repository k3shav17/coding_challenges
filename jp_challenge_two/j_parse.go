package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type TokenType int

const (
	LBRACE TokenType = iota
	RBRACE
	COLON
	COMMA
	EOF
	STRING
	WHITESPACE
	INVALID
)

type Token struct {
	Type    TokenType
	Literal string
}

type Lexer struct {
	input        string
	position     int  // position of the current char
	readPosition int  // point next to the current char
	ch           byte // current char
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) NextToken() Token {

	var token Token

	l.skipWhiteSpace()

	switch l.ch {
	case '{':
		token = newToken(LBRACE, l.ch)
	case '}':
		token = newToken(RBRACE, l.ch)
	case ':':
		token = newToken(COLON, l.ch)
	case '"':
		token.Type = STRING
		token.Literal = l.readString()
	case ',':
		token = newToken(COMMA, l.ch)
	default:
		token = newToken(INVALID, l.ch)
	}
	l.readChar()
	return token
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	if l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

type Parser struct {
	l *Lexer
}

func NewParse(l *Lexer) *Parser {
	return &Parser{l: l}
}

func (p *Parser) Parse() bool {
	tok := p.l.NextToken()
	if tok.Type != LBRACE {
		return false
	}

	for {
		tok = p.l.NextToken()
		if tok.Type == EOF {
			break
		}
		if tok.Type == RBRACE {
			return true
		}
	}
	return false

}

func main() {

	commandLineArgs := os.Args[:]
	filePath := commandLineArgs[1]

	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Errorf("Unable to read file %q", filePath)
	}

	input := string(bytes)
	lexer := NewLexer(input)
	parser := NewParse(lexer)

	if parser.Parse() {
		fmt.Println("valid JSON object")
	} else {
		fmt.Println("invalid JSON object")
	}
}
