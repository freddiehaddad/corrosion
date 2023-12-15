package parser

import (
	"fmt"
	"strconv"

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
	ASSIGN
	EQ
	LTGT
	SUM
	PRODUCT
	PREFIX
)

var precedences = map[token.TokenType]int{
	token.ASSIGN:   ASSIGN,
	token.LT:       LTGT,
	token.GT:       LTGT,
	token.LT_EQUAL: LTGT,
	token.GT_EQUAL: LTGT,
	token.EQ:       EQ,
	token.NOT_EQ:   EQ,
	token.MINUS:    SUM,
	token.PLUS:     SUM,
	token.MULTIPLY: PRODUCT,
	token.DIVIDE:   PRODUCT,
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
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MULTIPLY, p.parseInfixExpression)
	p.registerInfix(token.DIVIDE, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.LT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.GT_EQUAL, p.parseInfixExpression)

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseInteger)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

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
		return nil
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

func (p *Parser) parseAssignmentExpression(left ast.Expression) ast.Expression {
	a := &ast.AssignmentExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	if _, ok := left.(*ast.Identifier); !ok {
		msg := fmt.Sprintf("cannot assign to expression %q",
			left.String())
		p.error(msg)
		return nil
	}

	p.nextToken()
	a.Right = p.parseExpression(LOWEST)
	return a
}

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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // consume the (
	e := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentTokenIs(token.TRUE),
	}
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseInteger() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.error(err.Error())
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.currentToken,
		Value: value,
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
