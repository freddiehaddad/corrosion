package evaluator

import (
	"fmt"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/object"
)

var (
	NULL  = &object.Null{Value: nil}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type comparisonFunction func(string, object.Object, object.Object) object.Object

var comparisonFunctions map[object.ObjectType]comparisonFunction

func init() {
	comparisonFunctions = make(map[object.ObjectType]comparisonFunction)

	comparisonFunctions[object.BOOLEAN_OBJ] = compareBooleans
	comparisonFunctions[object.INTEGER_OBJ] = compareIntegers
}

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
	case *ast.AssignmentExpression:
		return evalAssignmentExpression(node, env)
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
	ie *ast.InfixExpression, env *object.Environment,
) object.Object {
	left := Eval(ie.Left, env)
	if checkEvalError(left) {
		return left
	}

	right := Eval(ie.Right, env)
	if checkEvalError(right) {
		return right
	}

	switch ie.Operator {
	case "+", "-", "*", "/":
		return evalArithmeticExpression(ie.Operator, left, right)
	case "==", "!=":
		return evalEqualityExpression(ie.Operator, left, right)
	case "<", "<=", ">", ">=":
		return evalRelationalExpression(ie.Operator, left, right)
	default:
		return evalError(fmt.Sprintf("ERROR: invalid operator=%q (%+v)",
			ie.Operator, ie))
	}
}

func evalAssignmentExpression(
	ae *ast.AssignmentExpression, env *object.Environment,
) object.Object {
	id, ok := ae.Left.(*ast.Identifier)
	if !ok {
		return evalError(
			fmt.Sprintf("runtime error. identifier expected (%+v)",
				ae.Left))
	}

	right := Eval(ae.Right, env)
	if checkEvalError(right) {
		return right
	}

	switch ae.Operator {
	case "=":
		env.Set(id.Value, right)
		return right
	default:
		return evalError(fmt.Sprintf("ERROR: invalid operator=%q (%+v)",
			ae.Operator, ae))
	}
}

// Checks of obj is an object.Integer type and returns the object.Integer form.
// Otherwise returns nil and an objectError as the second argument.
func expectIntegerObject(obj object.Object) (*object.Integer, object.Object) {
	i, ok := obj.(*object.Integer)
	if !ok {
		e := fmt.Sprintf("ERROR: expected integer, got=%T (%+v)",
			obj, obj)
		return nil, evalError(e)
	}
	return i, nil
}

func evalArithmeticExpression(
	op string, left, right object.Object,
) object.Object {
	value := &object.Integer{}

	l, err := expectIntegerObject(left)
	if err != nil {
		return err
	}

	r, err := expectIntegerObject(right)
	if err != nil {
		return err
	}

	switch op {
	case "+":
		value.Value = l.Value + r.Value
	case "-":
		value.Value = l.Value - r.Value
	case "*":
		value.Value = l.Value * r.Value
	case "/":
		if r.Value == 0 {
			return evalError(divisionByZero(left, right))
		}
		value.Value = l.Value / r.Value
	}

	return value
}

func mixedTypeError(op string, left, right object.Object) object.Object {
	e := fmt.Sprintf(`ERROR: comparison operation requires matching operand
		types. left=%s (%+v) %s right=%s (%+v)`,
		left.Type(), left, op, right.Type(), right)
	return &object.Error{Value: e}
}

func evalBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
}

func compareBooleans(op string, left, right object.Object) object.Object {
	l := left.(*object.Boolean)
	r := right.(*object.Boolean)

	var result bool

	switch op {
	case "==":
		result = l.Value == r.Value
	case "!=":
		result = l.Value != r.Value
	default:
		return evalError(
			fmt.Sprintf("unsupported comparison operator %s", op))
	}

	return evalBooleanObject(result)
}

func compareIntegers(op string, left, right object.Object) object.Object {
	l := left.(*object.Integer)
	r := right.(*object.Integer)

	var result bool

	switch op {
	case "==":
		result = l.Value == r.Value
	case "!=":
		result = l.Value != r.Value
	case "<":
		result = l.Value < r.Value
	case "<=":
		result = l.Value <= r.Value
	case ">":
		result = l.Value > r.Value
	case ">=":
		result = l.Value >= r.Value
	default:
		return evalError(
			fmt.Sprintf("unsupported comparison operator %s", op))
	}

	return evalBooleanObject(result)
}

func evalEqualityExpression(
	op string, left, right object.Object,
) object.Object {
	if left.Type() != right.Type() {
		return mixedTypeError(op, left, right)
	}

	fn, ok := comparisonFunctions[left.Type()]
	if !ok {
		return evalError(
			fmt.Sprintf("ERROR: no comparison function for %s",
				left.Type()))
	}

	return fn(op, left, right)
}

func evalRelationalExpression(
	op string, left, right object.Object,
) object.Object {
	if left.Type() != right.Type() {
		return mixedTypeError(op, left, right)
	}

	if left.Type() == object.BOOLEAN_OBJ {
		e := "ERROR: relationl comparison with boolean operands"
		return evalError(e)
	}

	fn, ok := comparisonFunctions[left.Type()]
	if !ok {
		return evalError(
			fmt.Sprintf("ERROR: no comparison function for %s",
				left.Type()))
	}

	return fn(op, left, right)
}

func divisionByZero(l, r object.Object) string {
	e := fmt.Sprintf("ERROR: divide by zero error in expression (%s / %s)",
		l.Inspect(), r.Inspect())
	return e
}

func evalPrefixExpression(
	pe *ast.PrefixExpression, env *object.Environment,
) object.Object {
	result := Eval(pe.Right, env)

	switch obj := result.(type) {
	case *object.Integer:
		if pe.Operator != "-" {
			e := fmt.Sprintf(
				"ERROR: unsupported operator=%q node=%T (%+v)",
				pe.Operator, result, result)
			return evalError(e)
		}
		return &object.Integer{Value: -obj.Value}
	case *object.Boolean:
		if pe.Operator != "!" {
			e := fmt.Sprintf(
				"ERROR: unsupported operator=%q node=%T (%+v)",
				pe.Operator, result, result)
			return evalError(e)
		}
		return &object.Boolean{Value: !obj.Value}
	default:
		return evalError(
			fmt.Sprintf("ERROR: unsupported node=%T (%+v)",
				result, result))
	}
}

func evalBooleanExpression(
	b *ast.Boolean, env *object.Environment,
) object.Object {
	return evalBooleanObject(b.Value)
}

func evalIntegerLiteral(
	i *ast.IntegerLiteral, env *object.Environment,
) object.Object {
	return &object.Integer{
		Value: i.Value,
	}
}

func evalIdentifier(
	i *ast.Identifier, env *object.Environment,
) object.Object {
	if obj, ok := env.Get(i.Value); ok {
		return obj
	}

	e := fmt.Sprintf("ERROR: undefined identifier=%q (%+v)", i.Value, i)
	return evalError(e)
}

func evalDeclarationStatement(
	node *ast.DeclarationStatement, env *object.Environment,
) object.Object {
	value := Eval(node.Value, env)
	if value.Type() == object.ERROR_OBJ {
		return value
	}

	if _, exists := env.Get(node.Name.Value); exists {
		e := fmt.Sprintf("ERROR: identifier=%q already defined.",
			node.Name.Value)
		return evalError(e)
	}
	env.Set(node.Name.Value, value)
	return NULL
}

func evalStatements(
	statements []ast.Statement, env *object.Environment,
) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)
	}

	return result
}

func evalError(s string) object.Object {
	e := object.Error{
		Value: s,
	}
	return &e
}
