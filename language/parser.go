package language

import (
	"fmt"
	"strconv"
)

type (
	prefixParseFn func() Node
	infixParseFn  func(Node) Node
)

const (
	_ int = iota
	LOWEST
	LOGICAL_OR
	LOGICAL_AND
	EQUALS
	LESSGREATER
	CONCATENTATION
	SUM
	PRODUCT
	SHELLPRECEDENCE
)

var precedences = map[TokenType]int{
	OR:          LOGICAL_OR,
	AND:         LOGICAL_AND,
	EQUAL:       EQUALS,
	NOTEQUAL:    EQUALS,
	LESSTHAN:    LESSGREATER,
	GREATERTHAN: LESSGREATER,
	LORETO:      LESSGREATER,
	GORETO:      LESSGREATER,
	CONCAT:      CONCATENTATION,
	PLUS:        SUM,
	MINUS:       SUM,
	SLASH:       PRODUCT,
	ASTERISK:    PRODUCT,
	MODULO:      PRODUCT,
}

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
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIdentifier() Node {
	return &Identifier{Value: p.currentToken.Literal}
}

func (p *Parser) parseStringLiteral() Node {
	return &StringLiteral{Value: p.currentToken.Literal}
}

func (p *Parser) parseExecStatement() Node {
	exec := &ExecStatement{}

	return exec
}

func (p *Parser) parseStatement() Node {
	switch p.currentToken.Type {
	case TASK:
		return p.parseTaskDefinition()
	case RUN:
		return p.parseExecStatement()
	case IF:
		return p.parseIfStatement()
	case FOREACH:
		return p.parseForEachStatement()
	case WHILE:
		return p.parseWhileStatement()
	case COMPILE:
		return p.parseCompileStatement()
	case IDENT:
		// Check if this is an assignment
		if p.peekTokenIs(ASSIGN) {
			return p.parseAssignmentStatement()
		}
		// Otherwise it's an expression statement
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseNumberLiteral() Node {
	lit := &NumberLiteral{}
	val, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("couldn't parser %q as float.", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = val

	return lit
}

func (p *Parser) parseInfixExpression(left Node) Node {
	expression := &BinaryOperation{
		Left:     left,
		Operator: p.currentToken.Literal,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseExpression(precedence int) Node {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(SEMICOLON) && !p.peekTokenIs(EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(tokenType TokenType) {
	msg := fmt.Sprintf("no prefix parse for %s found", tokenType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(token TokenType) bool {
	if p.peekToken.Type == token {
		return true
	}
	return false
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseTaskDefinition() *TaskDef {
	task := &TaskDef{}

	if !p.expectPeek(IDENT) {
		return nil
	}

	task.Name = p.currentToken.Literal

	// Check for dependencies
	if p.peekTokenIs(DEPENDENCY) {
		p.nextToken() // consume DEPENDENCY
		task.Dependencies = p.parseDependencies()
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}

	task.Body = p.parseBlockStatement()

	return task
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{}
	block.Statements = []Node{}

	p.nextToken()

	for !p.currentTokenIs(RBRACE) && p.currentTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}
