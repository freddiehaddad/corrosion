package parser

import (
	"fmt"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/lexer"
	"github.com/freddiehaddad/corrosion/pkg/token"
)

type Parser struct {
	l            *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}

	for !p.eof() {
		switch p.currentToken.Type {
		case token.INT:
			if stmt := p.parseDeclarationStatement(); stmt != nil {
				program.Statements = append(program.Statements, stmt)
			}
		case token.RETURN:
			if stmt := p.parseReturnStatement(); stmt != nil {
				program.Statements = append(program.Statements, stmt)
			}
		default:
			msg := fmt.Sprintf("ParseProgram: unexpected token=%s encountered", p.currentToken.Type)
			p.error(msg)
		}
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
// Error functions
// ----------------------------------------------------------------------------

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) error(msg string) {
	p.errors = append(p.errors, msg)
}

// ----------------------------------------------------------------------------
// Parsing functions
// ----------------------------------------------------------------------------

func (p *Parser) parseDeclarationStatement() ast.Statement {
	decType := p.currentToken // int

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	switch p.peekToken.Type {
	case token.ASSIGN:
		return p.parseVariableDeclarationStatement(decType)
	default:
		p.error(fmt.Sprintf("unexpected peekToken=%s", p.peekToken.Type))
		return nil
	}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	rs := &ast.ReturnStatement{Token: p.currentToken} // return

	// TODO: parse the expression.
	// skip everything up to the semicolon
	for !p.eof() {
		if p.currentTokenIs(token.SEMICOLON) {
			return rs
		}
		p.nextToken()
	}

	return rs
}

func (p *Parser) parseVariableDeclarationStatement(decType token.Token) ast.Statement {
	ds := &ast.DeclarationStatement{Token: decType}

	ds.Name = p.parseIdentifier()

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: parse the expression.
	// skip everything up to the semicolon
	for !p.eof() {
		if p.currentTokenIs(token.SEMICOLON) {
			return ds
		}
		p.nextToken()
	}

	p.error(fmt.Sprintf("parseDeclarationStatement: %+v - expected a semicolon after expression", ds))
	return nil
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}
