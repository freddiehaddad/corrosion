package evaluator

import (
	"fmt"

	"github.com/freddiehaddad/corrosion/pkg/ast"
	"github.com/freddiehaddad/corrosion/pkg/object"
)

// ----------------------------------------------------------------------------
// Constant objects
// ----------------------------------------------------------------------------

var (
	NULL  = &object.Null{Value: nil}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// ----------------------------------------------------------------------------
// Type comparisons help functions
// ----------------------------------------------------------------------------

type comparisonFunction func(string, object.Object, object.Object) object.Object

var comparisonFunctions map[object.ObjectType]comparisonFunction

func init() {
	comparisonFunctions = make(map[object.ObjectType]comparisonFunction)

	comparisonFunctions[object.BOOLEAN_OBJ] = compareBooleans
	comparisonFunctions[object.INTEGER_OBJ] = compareIntegers
}

// ----------------------------------------------------------------------------
// Evaluation entry point
// ----------------------------------------------------------------------------

// Evaluates the node and returns an object representing the expression value.
// Returns NULL object for non-value producing statements.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)
	case *ast.VariableDeclarationStatement:
		return evalDeclarationStatement(node, env)
	case *ast.FunctionDeclarationStatement:
		return evalFunctionDeclarationStatement(node, env)
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
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.FunctionCallExpression:
		return evalFunctionCallExpression(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	default:
		e := fmt.Sprintf("ERROR: unsupported node=%T (%+v)", node, node)
		return evalError(e)
	}
}

// ----------------------------------------------------------------------------
// Basic type evaluators
// ----------------------------------------------------------------------------

func evalBooleanExpression(
	b *ast.Boolean, env *object.Environment,
) object.Object {
	return evalBooleanObject(b.Value)
}

func evalBooleanObject(value bool) *object.Boolean {
	if value {
		return TRUE
	}
	return FALSE
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

// ----------------------------------------------------------------------------
// Expression evaluators
// ----------------------------------------------------------------------------

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
		obj, _ := env.Update(id.Value, right)
		return obj
	default:
		return evalError(fmt.Sprintf("ERROR: invalid operator=%q (%+v)",
			ae.Operator, ae))
	}
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

func evalFunctionCallExpression(
	node *ast.FunctionCallExpression, env *object.Environment,
) object.Object {
	function := Eval(node.Function, env)
	if checkEvalError(function) {
		return function
	}

	args := evalFunctionCallArguments(node.Arguments, env)
	if len(args) == 1 && checkEvalError(args[0]) {
		return args[0]
	}

	switch function := function.(type) {
	case *object.Function:
		extendedEnv := object.NewScopedEnvironment(function.Env)
		prepareFunctionCallParameters(
			args, function.Parameters, extendedEnv)

		evaluated := Eval(function.Body, extendedEnv)
		if ret, ok := evaluated.(*object.Return); ok {
			return ret.Value
		}
		return evaluated

	default:
		return evalError(fmt.Sprintf("not a function: %s",
			function.Type()))
	}
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

// ----------------------------------------------------------------------------
// Statement evaluators
// ----------------------------------------------------------------------------

func evalBlockStatement(
	node *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	for _, statement := range node.Statements {
		obj := Eval(statement, env)
		if obj.Type() == object.RETURN_OBJ {
			return obj
		}
	}
	return NULL
}

func evalDeclarationStatement(
	node *ast.VariableDeclarationStatement, env *object.Environment,
) object.Object {
	val := Eval(node.Value, env)
	if checkEvalError(val) {
		return val
	}

	if _, exists := env.Get(node.Name.Value); exists {
		e := fmt.Sprintf("ERROR: identifier=%q already defined.",
			node.Name.Value)
		return evalError(e)
	}
	env.Set(node.Name.Value, val)
	return NULL
}

func evalFunctionDeclarationStatement(
	node *ast.FunctionDeclarationStatement, env *object.Environment,
) object.Object {
	var function object.Function

	function.Parameters = node.Parameters
	function.Body = node.Body
	function.Env = env

	env.Set(node.Name.Value, &function)
	return NULL
}

func evalIfStatement(
	node *ast.IfStatement,
	env *object.Environment,
) object.Object {
	obj := Eval(node.Condition, env)

	condition, ok := obj.(*object.Boolean)
	if !ok {
		e := fmt.Sprintf(
			"ERROR: if condition must evaluate to a bool. got=%T",
			obj)
		return evalError(e)
	}

	if condition.Value {
		local := object.NewScopedEnvironment(env)
		return Eval(node.Consequence, local)
	}

	if node.Alternative != nil {
		local := object.NewScopedEnvironment(env)
		return Eval(node.Alternative, local)
	}

	return NULL
}

func evalReturnStatement(node *ast.ReturnStatement,
	env *object.Environment,
) object.Object {
	val := Eval(node.ReturnValue, env)
	if checkEvalError(val) {
		return val
	}

	return &object.Return{Value: val}
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

// ----------------------------------------------------------------------------
// Evaluator Helpers
// ----------------------------------------------------------------------------

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

// Evaluates the function call argument expressions and returns them in a slice
// to be set in the functions local environment scope.
func evalFunctionCallArguments(
	args []ast.Expression, env *object.Environment,
) []object.Object {
	obj := []object.Object{}

	for _, arg := range args {
		e := Eval(arg, env)
		if checkEvalError(e) {
			return []object.Object{e}
		}
		obj = append(obj, e)
	}

	return obj
}

// Copies function call arguments to the function's scope.
func prepareFunctionCallParameters(
	args []object.Object,
	params []ast.Identifier,
	env *object.Environment,
) {
	for index, param := range params {
		env.Set(param.Value, args[index])
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

// ----------------------------------------------------------------------------
// Error handling
// ----------------------------------------------------------------------------

func checkEvalError(obj object.Object) bool {
	switch obj.(type) {
	case *object.Error:
		return true
	}
	return false
}

func divisionByZero(l, r object.Object) string {
	e := fmt.Sprintf("ERROR: divide by zero error in expression (%s / %s)",
		l.Inspect(), r.Inspect())
	return e
}

func evalError(s string) object.Object {
	e := object.Error{
		Value: s,
	}
	return &e
}

func mixedTypeError(op string, left, right object.Object) object.Object {
	e := fmt.Sprintf(`ERROR: comparison operation requires matching operand
		types. left=%s (%+v) %s right=%s (%+v)`,
		left.Type(), left, op, right.Type(), right)
	return &object.Error{Value: e}
}
