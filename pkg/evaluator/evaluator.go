package evaluator

import (
	"fmt"
	"strconv"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/object"
)

var (
	NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	default:
		return NULL
	}
}

// ----------------------------------------------------------------------------
// Evaluators
// ----------------------------------------------------------------------------

func evalInfixExpression(ie *ast.InfixExpression) object.Object {
	result := &object.Integer{}

	left := Eval(ie.Left)
	right := Eval(ie.Right)

	lValue, _ := strconv.ParseInt(left.Inspect(), 10, 64)
	rValue, _ := strconv.ParseInt(right.Inspect(), 10, 64)

	var value string

	switch ie.Operator {
	case "+":
		v := lValue + rValue
		value = fmt.Sprintf("%d", v)
	case "-":
		v := lValue - rValue
		value = fmt.Sprintf("%d", v)
	default:
		return NULL
	}

	result.Value = value

	return result
}

func evalPrefixExpression(pe *ast.PrefixExpression) object.Object {
	result := Eval(pe.Right)

	switch result.(type) {
	case *object.Integer:
		if pe.Operator == "-" {
			s := result.Inspect()
			v, _ := strconv.ParseInt(s, 10, 64)
			v = -v
			s = fmt.Sprintf("%d", v)
			return &object.Integer{Value: s}
		}
	}

	return NULL
}

func evalIntegerLiteral(i *ast.IntegerLiteral) object.Object {
	return &object.Integer{Value: i.Value}
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)
	}

	return result
}
