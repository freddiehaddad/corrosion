package parser

import (
	"fmt"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/lexer"
	"github.com/freddiehaddad/corrosion/pkg/token"
)

// ----------------------------------------------------------------------------
// Operator precendene table
// ----------------------------------------------------------------------------

const (
	_ int = iota
	LOWEST
	SUM
	PRODUCT
	PREFIX
)

var precedences = map[token.TokenType]int{
	token.MINUS:    SUM,
	token.PLUS:     SUM,
	token.MULTIPLY: PRODUCT,
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

// ----------------------------------------------------------------------------
// Pratt parser semantic code
// ----------------------------------------------------------------------------

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// ----------------------------------------------------------------------------
// Parser
// ----------------------------------------------------------------------------

type Parser struct {
	l              *lexer.Lexer
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
	currentToken   token.Token
	peekToken      token.Token
	errors         []string
}

// ----------------------------------------------------------------------------
// Parser interface
// ----------------------------------------------------------------------------

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseInteger)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for !p.eof() {
		var stmt ast.Statement
		switch p.currentToken.Type {
		case token.INT:
			stmt = p.parseDeclarationStatement()
		case token.RETURN:
			stmt = p.parseReturnStatement()
		default:
			stmt = p.parseExpressionStatement()
		}
		program.Statements = append(program.Statements, stmt)
		p.nextToken()
	}

	return program
}

// ----------------------------------------------------------------------------
// Token functions
// ----------------------------------------------------------------------------

func (p *Parser) eof() bool {
	return p.currentTokenIs(token.EOF)
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)
	return false
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("peekToken=%s expected=%s", p.peekToken.Type, t)
	p.error(msg)
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ----------------------------------------------------------------------------
// Prefix functions
// ----------------------------------------------------------------------------

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(tt token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %q found", tt)
	p.error(msg)
}

// ----------------------------------------------------------------------------
// Error functions
// ----------------------------------------------------------------------------

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, msg)
}

// ----------------------------------------------------------------------------
// Statement parsing functions
// ----------------------------------------------------------------------------

func (p *Parser) parseDeclarationStatement() ast.Statement {
	var stmt ast.Statement
	decType := p.currentToken // int

	if !p.expectPeek(token.IDENT) {
		return stmt
	}

	switch p.peekToken.Type {
	case token.ASSIGN:
		stmt = p.parseVariableDeclarationStatement(decType)
	default:
		e := fmt.Sprintf("unexpected peekToken=%s", p.peekToken.Type)
		p.error(e)
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	rs := &ast.ReturnStatement{Token: p.currentToken} // return

	p.nextToken()

	rs.ReturnValue = p.parseExpression(LOWEST)

	p.expectPeek(token.SEMICOLON)

	return rs
}

func (p *Parser) parseVariableDeclarationStatement(
	decType token.Token,
) ast.Statement {
	ds := &ast.DeclarationStatement{Token: decType} // int

	ds.Name = ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	} // x

	if !p.expectPeek(token.ASSIGN) {
		return ds
	}
	p.nextToken() // =

	ds.Value = p.parseExpression(LOWEST)

	p.expectPeek(token.SEMICOLON)

	return ds
}

// ----------------------------------------------------------------------------
// Expression parsing functions
// ----------------------------------------------------------------------------

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	e := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	e.Right = p.parseExpression(precedence)
	return e
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.nextToken() // -
	pe.Right = p.parseExpression(PREFIX)
	return pe
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	return &ast.IntegerLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	es := &ast.ExpressionStatement{Token: p.currentToken}
	es.Expression = p.parseExpression(LOWEST)
	p.expectPeek(token.SEMICOLON)
	return es
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	// 10 + 10; precedence: LOWEST
	prefix := p.prefixParseFns[p.currentToken.Type] // INTEGER
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix() // parseInteger() -> 10

	// peekToken: +; LOWEST < peekPrecedence
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		// parseInfixExpression()
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken() // currentToken: +
		// parseInfixExpression(10) ->
		//   InfixExpression{T: 10, L: 10 O: + R: 10}
		leftExp = infix(leftExp)
	}

	return leftExp
}
