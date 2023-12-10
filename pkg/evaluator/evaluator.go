package evaluator

import (
	"fmt"
	"strconv"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/object"
)

var NULL = &object.Null{}

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
	case *ast.Boolean:
		return evalBooleanExpression(node, env)
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

func checkEvalError(obj object.Object) bool {
	switch obj.(type) {
	case *object.Error:
		return true
	}
	return false
}

func evalInfixExpression(
	ie *ast.InfixExpression,
	env *object.Environment,
) object.Object {
	result := &object.Integer{}

	left := Eval(ie.Left, env)
	if checkEvalError(left) {
		return left
	}

	right := Eval(ie.Right, env)
	if checkEvalError(right) {
		return right
	}

	lValue, err := strconv.ParseInt(left.Inspect(), 10, 64)
	if err != nil {
		return evalError(err.Error())
	}

	rValue, err := strconv.ParseInt(right.Inspect(), 10, 64)
	if err != nil {
		return evalError(err.Error())
	}

	var value string

	switch ie.Operator {
	case "+":
		v := lValue + rValue
		value = fmt.Sprintf("%d", v)
	case "-":
		v := lValue - rValue
		value = fmt.Sprintf("%d", v)
	case "*":
		v := lValue * rValue
		value = fmt.Sprintf("%d", v)
	case "/":
		if rValue == 0 {
			return evalError(divisionByZero(ie))
		}
		v := lValue / rValue
		value = fmt.Sprintf("%d", v)
	default:
		return evalError(fmt.Sprintf("ERROR: invalid operator=%q (%+v)",
			ie.Operator, ie))
	}

	result.Value = value

	return result
}

func divisionByZero(ie *ast.InfixExpression) string {
	e := fmt.Sprintf("ERROR: divide by zero error in expression (%+v)", ie)
	return e
}

func evalPrefixExpression(
	pe *ast.PrefixExpression,
	env *object.Environment,
) object.Object {
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
	case *object.Boolean:
		if pe.Operator != "!" {
			e := fmt.Sprintf(
				"ERROR: unsupported operator=%q node=%T (%+v)",
				pe.Operator, result, result)
			return evalError(e)
		}
		s := result.Inspect()
		v, _ := strconv.ParseBool(s)
		v = !v
		s = fmt.Sprintf("%t", v)
		return &object.Boolean{Value: s}
	default:
		return evalError(fmt.Sprintf("ERROR: unsupported node=%T (%+v)",
			result, result))
	}
}

func evalBooleanExpression(
	b *ast.Boolean,
	env *object.Environment,
) object.Object {
	return &object.Boolean{Value: b.Value}
}

func evalIntegerLiteral(
	i *ast.IntegerLiteral,
	env *object.Environment,
) object.Object {
	return &object.Integer{Value: i.Value}
}

func evalIdentifier(i *ast.Identifier, env *object.Environment) object.Object {
	if obj, ok := env.Get(i.Value); ok {
		return obj
	}

	e := fmt.Sprintf("ERROR: undefined identifier=%q (%+v)", i.Value, i)
	return evalError(e)
}

func evalDeclarationStatement(
	node *ast.DeclarationStatement,
	env *object.Environment,
) object.Object {
	value := Eval(node.Value, env)
	env.Set(node.Name.Value, value)
	return NULL
}

func evalStatements(
	statements []ast.Statement,
	env *object.Environment,
) object.Object {
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
