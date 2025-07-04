package language

import (
	"fmt"
	"strings"
)

type (
	prefixParseFn func() Node
	infixParseFn  func(Node) Node
)

type Parser struct {
	l            *Lexer
	currentToken Token
	peekToken    Token
	errors       []string

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []string{},
		prefixParseFns: make(map[TokenType]prefixParseFn),
		infixParseFns:  make(map[TokenType]infixParseFn),
	}

	// Register prefix parse functions
	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(NUMBER, p.parseNumberLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(NOT, p.parsePrefixExpression)
	p.registerPrefix(SHELL, p.parseShellExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)

	// Register infix parse functions
	p.registerInfix(PLUS, p.parseInfixExpression)
	p.registerInfix(MINUS, p.parseInfixExpression)
	p.registerInfix(SLASH, p.parseInfixExpression)
	p.registerInfix(ASTERISK, p.parseInfixExpression)
	p.registerInfix(MODULO, p.parseInfixExpression)
	p.registerInfix(EQUAL, p.parseInfixExpression)
	p.registerInfix(NOTEQUAL, p.parseInfixExpression)
	p.registerInfix(LESSTHAN, p.parseInfixExpression)
	p.registerInfix(GREATERTHAN, p.parseInfixExpression)
	p.registerInfix(LORETO, p.parseInfixExpression)
	p.registerInfix(GORETO, p.parseInfixExpression)
	p.registerInfix(AND, p.parseInfixExpression)
	p.registerInfix(OR, p.parseInfixExpression)
	p.registerInfix(CONCAT, p.parseConcatExpression)

	// Initialize peek and current token
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// Skip newlines and comments
	for p.peekToken.Type == NEWLINE || p.peekToken.Type == COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) Errors() []string {
	return p.errors // strings holding human readable errors (With a lot of newlines yk)
}

func (p *Parser) peekError(t TokenType) {
	var out strings.Builder
	out.WriteString("\n")
	out.WriteString("expected %s but got %s\n")
	out.WriteString(fmt.Sprintf("\t->%s:%d:%d :\n", p.l.filename, p.l.line, p.l.column))
	p.errors = append(p.errors, out.String())
}
