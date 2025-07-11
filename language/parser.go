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

	p.registerPrefix(IDENT, p.parseIdentifier)
	p.registerPrefix(STRING, p.parseStringLiteral)
	p.registerPrefix(NUMBER, p.parseNumberLiteral)
	p.registerPrefix(MINUS, p.parsePrefixExpression)
	p.registerPrefix(NOT, p.parsePrefixExpression)
	p.registerPrefix(SHELL, p.parseShellExpression)
	p.registerPrefix(LPAREN, p.parseGroupedExpression)

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

	p.nextToken()
	p.nextToken()

	return p
}

// -- Registration --
func (p *Parser) registerPrefix(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// -- Utility --
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
	for p.peekToken.Type == NEWLINE || p.peekToken.Type == COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) peekTokenIs(token TokenType) bool {
	return p.peekToken.Type == token
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) noPrefixParseFnError(tokenType TokenType) {
	msg := fmt.Sprintf("no prefix parse for %s found", tokenType)
	p.errors = append(p.errors, msg)
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

// -- Expressions --
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

func (p *Parser) parsePrefixExpression() Node {
	expression := &UnaryOperation{Operator: p.currentToken.Literal}
	p.nextToken()
	expression.Operand = p.parseExpression(LOWEST)
	return expression
}

func (p *Parser) parseInfixExpression(left Node) Node {
	expression := &BinaryOperation{Left: left, Operator: p.currentToken.Literal}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() Node {
	p.nextToken()
	expr := p.parseExpression(LOWEST)
	if !p.expectPeek(RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseConcatExpression(left Node) Node {
	expr := &ConcatOperation{Left: left}
	p.nextToken()
	expr.Right = p.parseExpression(CONCATENTATION)
	return expr
}

func (p *Parser) parseIdentifier() Node {
	return &Identifier{Value: p.currentToken.Literal}
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

func (p *Parser) parseStringLiteral() Node {
	return &StringLiteral{Value: p.currentToken.Literal}
}

func (p *Parser) parseShellExpression() Node {
	if !p.expectPeek(IDENT) {
		return nil
	}
	return &ShellExpr{Name: p.currentToken.Literal}
}

// -- Statements --
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
		if p.peekTokenIs(ASSIGN) {
			return p.parseAssignmentStatement()
		}
		return p.parseExpressionStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() Node {
	return p.parseExpression(LOWEST)
}

func (p *Parser) parseAssignmentStatement() *AssignmentStatement {
	assign := &AssignmentStatement{Name: p.currentToken.Literal}
	if !p.expectPeek(ASSIGN) {
		return nil
	}
	p.nextToken()
	assign.Value = p.parseExpression(LOWEST)
	return assign
}

func (p *Parser) parseExecStatement() *ExecStatement {
	exec := &ExecStatement{}
	if !p.expectPeek(IDENT) {
		return nil
	}
	exec.TaskName = p.currentToken.Literal
	return exec
}

func (p *Parser) parseIfStatement() *IfStatement {
	ifStmt := &IfStatement{}
	p.nextToken()
	ifStmt.Condition = p.parseExpression(LOWEST)

	if p.peekTokenIs(ELSE) {
		p.nextToken()
		if !p.expectPeek(LBRACE) {
			return nil
		}
		ifStmt.ElseBlock = p.parseBlockStatement()
	}
	return ifStmt
}

func (p *Parser) parseForEachStatement() *ForEachStatement {
	forEach := &ForEachStatement{}
	if !p.expectPeek(IDENT) {
		return nil
	}
	forEach.VarName = p.currentToken.Literal
	if !p.expectPeek(STRING) && !p.expectPeek(IDENT) {
		return nil
	}
	forEach.Pattern = p.currentToken.Literal
	if !p.expectPeek(LBRACE) {
		return nil
	}
	forEach.Body = p.parseBlockStatement()
	return forEach
}

func (p *Parser) parseWhileStatement() Node {
	return &Identifier{Value: "PLACEHOLDER... IN IMPL"}
}

func (p *Parser) parseCompileStatement() *CompileStatement {
	cmpStmt := &CompileStatement{}
	p.nextToken()
	cmpStmt.File = p.parseExpression(LOWEST)
	p.nextToken()
	cmpStmt.Command = p.parseExpression(LOWEST)
	return cmpStmt
}

// -- Task Definition --
func (p *Parser) parseTaskDefinition() *TaskDef {
	task := &TaskDef{}
	if !p.expectPeek(IDENT) {
		return nil
	}
	task.Name = p.currentToken.Literal
	if p.peekTokenIs(DEPENDENCY) {
		p.nextToken()
		p.parseDependencies(task)
	}
	if !p.expectPeek(LBRACE) {
		return nil
	}
	task.Body = p.parseBlockStatement()
	return task
}

func (p *Parser) parseDependencies(t *TaskDef) {
	if !p.expectPeek(IDENT) {
		return
	}
	t.Dependencies = append(t.Dependencies, p.currentToken.Literal)
	for p.peekTokenIs(COMMA) {
		p.nextToken()
		if !p.expectPeek(COMMA) {
			break
		}
		t.Dependencies = append(t.Dependencies, p.currentToken.Literal)
	}
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{Statements: []Node{}}
	p.nextToken()
	for !p.currentTokenIs(RBRACE) && !p.currentTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}
