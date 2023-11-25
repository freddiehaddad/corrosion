package ast

import (
	"strings"

	"github.com/freddiehaddad/corrosion/pkg/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	sb := strings.Builder{}

	for _, s := range p.Statements {
		sb.WriteString(s.String())
	}

	return sb.String()
}

// ----------------------------------------------------------------------------
// Statement types
// ----------------------------------------------------------------------------

type DeclarationStatement struct {
	Value Expression
	Name  *Identifier
	Token token.Token
}

func (ds *DeclarationStatement) statementNode() {}
func (ds *DeclarationStatement) TokenLiteral() string {
	return ds.Token.Literal
}
func (ds *DeclarationStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString(ds.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(ds.Name.String())
	sb.WriteString(" = ")
	if ds.Value != nil {
		sb.WriteString(ds.Value.String())
	}
	sb.WriteString(";")
	return sb.String()
}

type ExpressionStatement struct {
	Expression Expression
	Token      token.Token
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

type ReturnStatement struct {
	ReturnValue Expression
	Token       token.Token
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
func (rs *ReturnStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString(rs.TokenLiteral())
	sb.WriteString(" ")
	if rs.ReturnValue != nil {
		sb.WriteString(rs.ReturnValue.String())
	}
	sb.WriteString(";")
	return sb.String()
}

// ----------------------------------------------------------------------------
// Basic types
// ----------------------------------------------------------------------------

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value string
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string {
	return i.Value
}
