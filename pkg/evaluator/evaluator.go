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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.DeclarationStatement:
		return evalDeclarationStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	default:
		e := fmt.Sprintf("ERROR: unsupported node=%T (%+v)", node, node)
		return evalError(e)
	}
}

// ----------------------------------------------------------------------------
// Evaluators
// ----------------------------------------------------------------------------

func evalInfixExpression(
	ie *ast.InfixExpression,
	env *object.Environment) object.Object {

	result := &object.Integer{}

	left := Eval(ie.Left, env)
	right := Eval(ie.Right, env)

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
		return evalError(fmt.Sprintf("ERROR: invalid operator=%q (%+v)",
			ie.Operator, ie))
	}

	result.Value = value

	return result
}

func evalPrefixExpression(
	pe *ast.PrefixExpression,
	env *object.Environment) object.Object {

	result := Eval(pe.Right, env)

	switch result.(type) {
	case *object.Integer:
		if pe.Operator == "-" {
			s := result.Inspect()
			v, _ := strconv.ParseInt(s, 10, 64)
			v = -v
			s = fmt.Sprintf("%d", v)
			return &object.Integer{Value: s}
		} else {
			e := fmt.Sprintf(
				"ERROR: unsupported operator=%q node=%T (%+v)",
				pe.Operator, result, result)
			return evalError(e)
		}
	default:
		return evalError(fmt.Sprintf("ERROR: unsupported node=%T (%+v)",
			result, result))
	}
}

func evalIntegerLiteral(
	i *ast.IntegerLiteral,
	env *object.Environment) object.Object {

	return &object.Integer{Value: i.Value}
}

func evalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	obj, _ := env.Get(i.Value)
	return obj
}

func evalDeclarationStatement(
	node *ast.DeclarationStatement,
	env *object.Environment) object.Object {

	value := Eval(node.Value, env)
	env.Set(node.Name.Value, value)
	return NULL
}

func evalStatements(
	statements []ast.Statement,
	env *object.Environment) object.Object {

	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)
	}

	return result
}

func evalError(s string) object.Object {
	e := object.Error{Value: s}
	return &e
}
