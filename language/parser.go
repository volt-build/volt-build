package language

import (
	"fmt"
	"strconv"
)

type Parser struct {
	l            *Lexer
	currentToken Token
	peekToken    Token
	errors       []string
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// To init peek and current token.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()

	// no newlines and comments (skipping here)
	for p.peekToken.Type == NEWLINE || p.peekToken.Type == COMMENT {
		p.peekToken = p.l.NextToken()
	}
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
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
	msg := fmt.Sprintf("Line %d, Column %d: expected next token to be %v, got %v instead", p.peekToken.Line, p.peekToken.Column, lexerMap[t], lexerMap[p.peekToken.Type])
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{
		Statements: []Node{},
	}

	for !p.currentTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Node {
	switch p.currentToken.Type {
	case TASK:
		return p.parseTaskDefinition()
	case IDENT:
		if p.peekTokenIs(ASSIGN) {
			return p.parseAssignStatement()
		}

		if p.currentToken.Literal == "exec" {
			return p.parseExecStatement()
		}
		if p.currentToken.Literal == "shell" {
			return p.parseShellStatement()
		}
		if p.currentToken.Literal == "push" {
			return p.parsePushStatement()
		}
		if p.currentToken.Literal == "compile" {
			return p.parseCompileStatement()
		}
		return nil
	case COMPILE:
		return p.parseCompileStatement()
	case IF:
		return p.parseIfStatement()
	case FOREACH:
		return p.parseForEachStatement()
	case SHELL:
		return p.parseShellStatement()

	default:
		return nil
	}
}

func (p *Parser) parseCompileStatement() *CompileStatement {
	stmt := &CompileStatement{}

	p.nextToken()                   // consume the "compile" keyword"
	stmt.File = p.parseExpression() // parse the expression in front of `compile`
	p.nextToken()                   // move onto the next expression

	stmt.Command = p.parseExpressionWithConcat()
	return stmt
}

func (p *Parser) parseExpressionWithConcat() Node {
	left := p.parseExpression()
	if p.peekTokenIs(CONCAT) {
		p.nextToken() // consume current token
		p.nextToken() // consume "compile"

		right := p.parseExpressionWithConcat()

		return &ConcatOperation{
			Left:  left,
			Right: right,
		}
	}

	return left
}

func (p *Parser) parseTaskDefinition() *TaskDef {
	task := &TaskDef{}
	task.Dependencies = []string{} // init a empty slice for now.

	if !p.expectPeek(IDENT) {
		return nil
	}

	task.Name = p.currentToken.Literal

	if p.peekTokenIs(IDENT) && p.peekToken.Literal == "requires" {
		p.nextToken() // consume "requires"

		if !p.expectPeek(IDENT) {
			return nil
		}

		task.Dependencies = append(task.Dependencies, p.currentToken.Literal)

		for p.peekTokenIs(COMMA) {
			p.nextToken()
			if !p.expectPeek(IDENT) {
				return nil
			}
			task.Dependencies = append(task.Dependencies, p.currentToken.Literal)
		}
	}
	if !p.expectPeek(LBRACE) {
		return nil
	}

	task.Body = p.parseBlockStatement()
	return task
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	block := &BlockStatement{
		Statements: []Node{},
	}

	p.nextToken() // comsume `{`

	for !p.currentTokenIs(RBRACE) && !p.currentTokenIs(EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExecStatement() *ExecStatement {
	stmt := &ExecStatement{}
	if !p.expectPeek(IDENT) {
		return nil
	}
	stmt.TaskName = p.currentToken.Literal
	return stmt
}

func (p *Parser) parseShellStatement() *ShellStatement {
	stmt := &ShellStatement{}
	p.nextToken()
	stmt.Command = p.parseExpression()

	return stmt
}

func (p *Parser) parsePushStatement() *PushStatement {
	stmt := &PushStatement{}
	p.nextToken()
	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseIfStatement() *IfStatement {
	stmt := &IfStatement{}

	p.nextToken() // consume if
	stmt.Condition = p.parseExpression()

	if !p.expectPeek(LBRACE) {
		return nil
	}

	stmt.ThenBlock = p.parseBlockStatement()

	if p.peekTokenIs(ELSE) {
		p.nextToken() // consume else

		if !p.expectPeek(LBRACE) {
			return nil
		}
		stmt.ElseBlock = p.parseBlockStatement()
	}

	return stmt
}

func (p *Parser) parseForEachStatement() *ForEachStatement {
	stmt := &ForEachStatement{}

	p.nextToken() // consume `foreach`

	if !p.currentTokenIs(STRING) && !p.currentTokenIs(IDENT) {
		p.errors = append(p.errors, fmt.Sprintf("Line %d, Column %d: expected string or identifier, got %s", p.currentToken.Line, p.currentToken.Column, p.currentToken.Literal))
		return nil
	}

	stmt.Pattern = p.currentToken.Literal
	stmt.VarName = "it"

	if p.peekTokenIs(IDENT) {
		p.nextToken()
		varName := p.currentToken.Literal
		if string(varName[0]) == "$" {
			varName = varName[1:]
		}
		stmt.VarName = varName
	}

	if !p.expectPeek(LBRACE) {
		return nil
	}
	stmt.Body = p.parseBlockStatement()
	return stmt
}

func (p *Parser) parseAssignStatement() *AssignmentStatement {
	stmt := &AssignmentStatement{
		Name: p.currentToken.Literal,
	}

	p.nextToken()
	p.nextToken()

	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseExpression() Node {
	var left Node

	switch p.currentToken.Type {
	case STRING:
		left = &StringLiteral{Value: p.currentToken.Literal}
	case NUMBER:
		value, _ := strconv.ParseFloat(p.currentToken.Literal, 64)
		left = &NumberLiteral{Value: value}
	case IDENT:
		left = &Identifier{Value: p.currentToken.Literal}
	case SHELL:
		p.nextToken() // consume '$'
		left = &ShellExpr{Name: p.currentToken.Literal}
	default:
		return nil
	}

	// Check for concatenation operator
	if p.peekTokenIs(CONCAT) {
		p.nextToken() // consume the current token
		p.nextToken() // consume the CONCAT token

		right := p.parseExpression() // Handle right-associative concatenation

		return &ConcatOperation{
			Left:  left,
			Right: right,
		}
	}

	return left
}
