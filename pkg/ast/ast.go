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
	Name  Identifier
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
	sb.WriteString(ds.Value.String())
	sb.WriteString(";")
	return sb.String()
}

type ExpressionStatement struct {
	Expression Expression
	Token      token.Token
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}
func (es *ExpressionStatement) String() string { return es.Expression.String() }

type AssignmentExpression struct {
	Token    token.Token
	Left     Expression
	Right    Expression
	Operator string
}

func (a *AssignmentExpression) expressionNode()      {}
func (a *AssignmentExpression) TokenLiteral() string { return a.Token.Literal }
func (a *AssignmentExpression) String() string {
	sb := strings.Builder{}
	sb.WriteByte('(')
	sb.WriteString(a.Left.String())
	sb.WriteByte(' ')
	sb.WriteString(a.Operator)
	sb.WriteByte(' ')
	sb.WriteString(a.Right.String())
	sb.WriteByte(')')
	return sb.String()
}

type InfixExpression struct {
	Left     Expression
	Right    Expression
	Token    token.Token
	Operator string
}

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }
func (i *InfixExpression) String() string {
	sb := strings.Builder{}
	sb.WriteByte('(')
	sb.WriteString(i.Left.String())
	sb.WriteByte(' ')
	sb.WriteString(i.Operator)
	sb.WriteByte(' ')
	sb.WriteString(i.Right.String())
	sb.WriteByte(')')
	return sb.String()
}

type PrefixExpression struct {
	Right    Expression
	Token    token.Token
	Operator string
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	sb := strings.Builder{}
	sb.WriteByte('(')
	sb.WriteString(p.Operator)
	sb.WriteString(p.Right.String())
	sb.WriteByte(')')
	return sb.String()
}

type ReturnStatement struct {
	ReturnValue Expression
	Token       token.Token
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString(rs.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(rs.ReturnValue.String())
	sb.WriteString(";")
	return sb.String()
}

type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
	Token       token.Token
}

func (is *IfStatement) statementNode()       {}
func (is *IfStatement) expressionNode()      {}
func (is *IfStatement) TokenLiteral() string { return is.Token.Literal }
func (is *IfStatement) String() string {
	var sb strings.Builder
	sb.WriteString("if")
	sb.WriteString(is.Condition.String())
	sb.WriteString(" ")
	sb.WriteString(is.Consequence.String())
	if is.Alternative != nil {
		sb.WriteString("else ")
		sb.WriteString(is.Alternative.String())
	}
	return sb.String()
}

type BlockStatement struct {
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var sb strings.Builder
	for _, s := range bs.Statements {
		sb.WriteString(s.String())
	}
	return sb.String()
}

// ----------------------------------------------------------------------------
// Basic types
// ----------------------------------------------------------------------------

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }
