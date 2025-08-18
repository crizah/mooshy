package evaluator

// ERROR HANDELING NEED TO BE FIXED

import (
	"mooshy/ast"
	"mooshy/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
	WRONG = &object.Error{Msg: "Cant use BuiltIn functions as Identifiers"}
)

func evalNotOp(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	default: // for !5
		return FALSE
	}
}

func evalMinusOp(right object.Object) object.Object {
	i, ok := right.(*object.Integer)
	if !ok {
		return &object.Error{Msg: "not integer object"}
		// return NULL
	}
	return &object.Integer{Value: -i.Value}

}

func evaluateOp(right object.Object, operator string) object.Object {
	switch operator {
	case "!":
		return evalNotOp(right)
	case "-":
		return evalMinusOp(right)
	default:
		return &object.Error{Msg: "not a prefix operator. got :" + operator + ". ! and - supported."}
		// return NULL

	}

}

func evalInfix(r object.Object, l object.Object, operator string) object.Object {

	switch operator {
	case "+":
		return evalInfixOp(l, r, operator)
	case "-":
		return evalInfixOp(l, r, operator)
	case "*":
		return evalInfixOp(l, r, operator)
	case "/":
		return evalInfixOp(l, r, operator)
	case "<":
		return evalInfixOp(l, r, operator)
	case ">":
		return evalInfixOp(l, r, operator)
	case "==":
		return evalInfixOp(l, r, operator)
	case "!=":
		return evalInfixOp(l, r, operator)

	default:
		return &object.Error{Msg: "unknown operator: " + operator}
		// return NULL
	}

}

func evalPostfix(left object.Object, operator string) object.Object {

	l, ok := left.(*object.Integer)
	if !ok {
		return &object.Error{Msg: "Cant perform postfix operatios on non integer objects"}
	}

	switch operator {
	case "++":
		return &object.Integer{Value: l.Value + 1}
	case "--":
		return &object.Integer{Value: l.Value - 1}
	default:
		return &object.Error{Msg: "Not a valid postfix operator"}

	}
}

func evalInfixOp(left object.Object, right object.Object, operator string) object.Object {
	// REFINE THIS WTF EVEN IS THIS CODE.

	// both left and right dont NEED rto be integres. can be of string type and return will be concacted string.
	// but for now, keep them both as only integers allowed

	if right.Type() != left.Type() {
		return &object.Error{Msg: "Type mismatch: " + right.Inspect() + " and " + left.Inspect()}
	}

	if right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ {
		// also do if of differetent object types and == and != operators instead of null

		r, ok := right.(*object.Integer)
		if !ok {
			return &object.Error{Msg: "Not Integer Object"}

			// return NULL
		}

		l, ok := left.(*object.Integer)
		if !ok {
			return &object.Error{Msg: "Not Integer Object"}
			// return NULL
		}
		switch operator {
		case "+":
			return &object.Integer{Value: l.Value + r.Value}
		case "-":
			return &object.Integer{Value: l.Value - r.Value}
		case "*":
			return &object.Integer{Value: l.Value * r.Value}
		case "/":
			return &object.Integer{Value: l.Value / r.Value}

		case "<":
			if l.Value < r.Value {
				return TRUE
			}
			return FALSE
		case ">":
			if l.Value > r.Value {
				return TRUE
			}
			return FALSE
		case "==":
			if l.Value == r.Value {
				return TRUE
			}
			return FALSE
		case "!=":
			if l.Value != r.Value {
				return TRUE
			}
			return FALSE
		default:
			return &object.Error{Msg: "Unknown operator: " + operator}
			// return NULL
		}
	} else if right.Type() == object.BOOL_OBJ && left.Type() == object.BOOL_OBJ {

		r, ok := right.(*object.Bool)
		if !ok {
			return &object.Error{Msg: "Not Bool Object"}
			// return NULL
		}

		l, ok := left.(*object.Bool)
		if !ok {
			return &object.Error{Msg: "Not Bool Object"}
			// return NULL
		}
		switch operator {
		case "==":
			if l.Value == r.Value {
				return TRUE
			}
			return FALSE
		case "!=":
			if l.Value != r.Value {
				return TRUE
			}
			return FALSE
		default:
			return &object.Error{Msg: "Cant perform " + operator + " on Bool objects"}
		}

	} else if right.Type() == object.STRING_OBJ && left.Type() == object.STRING_OBJ {
		r, ok := right.(*object.String)
		if !ok {
			return &object.Error{Msg: "Not String Object"}

			// return NULL
		}

		l, ok := left.(*object.String)
		if !ok {
			return &object.Error{Msg: "Not String Object"}
			// return NULL
		}
		switch operator {
		case "+":
			return &object.String{Value: l.Value + r.Value}

		case "==":
			if l.Value == r.Value {
				return TRUE
			}
			return FALSE
		case "!=":
			if l.Value != r.Value {
				return TRUE
			}
			return FALSE
		default:
			return &object.Error{Msg: "Cant perform " + operator + " on String objects"}
		}
	}
	// return NULL

	return &object.Error{Msg: "Unrecogniused Object Type "}

	// return NULL

}

func evalIfExpressions(ie *ast.IfExpression, env *object.Enviorment) object.Object {
	// 	if (x) {
	// puts("everything okay!");
	// } else {
	// puts("x is too high!");
	// shutdownSystem();
	// }

	condition := Eval(ie.Condition, env)
	if isTrue(condition) {
		return Eval(ie.Consequence, env)
	} else {
		if ie.Alternative != nil {
			return Eval(ie.Alternative, env)
		}
		return NULL
	}
}

func isTrue(condition object.Object) bool {
	switch condition {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}
func Eval(node ast.Node, env *object.Enviorment) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BoolExpression:
		return helper(node.Value) // so that we dont have to make a new instance everytime
		// return &object.Bool{Value: node.Value}

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
		// return evalStatements(node.Statements)
	case *ast.LetStatement:
		switch node.Name.Value {
		case "len":
			return WRONG
		case "sum":
			return WRONG
		case "append":
			return WRONG

		default:
			val := Eval(node.Value, env)
			return env.Put(node.Name.Value, val)

		}
	case *ast.ReAssignExpression:
		val := Eval(node.Value, env)
		id, ok := node.Name.(*ast.Identifier)
		if ok {

			_, c := env.Get(id.String())
			if !c {
				return &object.Error{Msg: "Cant assign an undeclared value"}

			}

			return env.Put(id.String(), val)

			// if env.Get(node.Name.String()) == nil, false {
			// 	return &object.Error{Msg: "Cant assign an undeclared value"}

			// }

		}
		return &object.Error{Msg: "Not an identifier" + node.Name.String()}

		// return env.Put(node.Name.String(), val)

	case *ast.FunctionLiteral:
		return &object.Function{Params: node.Parameter, Body: node.Body, Env: env}

	case *ast.CallExpression:

		fun := Eval(node.Function, env) // of this is len, then trigger builtIn
		args := evalArguments(node.Arguments, env)
		return evalOuterFunc(fun, args)

	case *ast.IndexExpression:
		arr := Eval(node.Name, env)
		// needs to return object of the same typeas arrs first one. we will confirm only same type allowed
		if ar_obj, ok := arr.(*object.Array); ok {
			if node.Index >= int64(len(ar_obj.Value)) || node.Index < 0 {
				return &object.Error{Msg: "Index out of bounds"}

			}
			idx := ar_obj.Value[node.Index]
			return idx
		}

		return &object.Error{Msg: "Not an Array Object"}

	case *ast.Identifier:
		return evalIdentifier(node, env)
		// can also be builInFunctions

	case *ast.IfExpression:
		return evalIfExpressions(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		return &object.Return{Value: val}

	case *ast.ArrayExpression:
		args := evalArguments(node.Value, env)

		for _, arg := range args {
			if args[0].Type() != arg.Type() {
				return &object.Error{Msg: "Different Types Array not Allowed"}
			}

		}

		if args[0].Type() == object.STRING_OBJ || args[0].Type() == object.INTEGER_OBJ {
			return &object.Array{Value: args}
		}

		return &object.Error{Msg: "Objects not of type String or Integer"}

	case *ast.PrefixExpression: // can be - or !
		right := Eval(node.Right, env)
		return evaluateOp(right, node.Operator)
	case *ast.InfixExpression:
		right := Eval(node.Right, env)
		left := Eval(node.Left, env)
		return evalInfix(right, left, node.Operator)
	case *ast.PostfixExpression:
		left := Eval(node.Left, env)
		return evalPostfix(left, node.Operator)
	case *ast.Program:
		// return evalProgram(node)
		return evalStatements(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	}

	return NULL
}

func evalIdentifier(
	node *ast.Identifier,
	env *object.Enviorment,
) object.Object {

	if val, ok := env.Get(node.Value); ok {
		return val
	} else if b, ok := builtins[node.Value]; ok {
		return b // this is a function object. we need to call Fn
	}

	return &object.Error{Msg: "identifier not found: " + node.Value}

}

func evalOuterFunc(obj object.Object, args []object.Object) object.Object {

	switch fn := obj.(type) {
	case *object.Function:
		extended := extendEnv(fn, args)
		evaluated := Eval(fn.Body, extended)
		return unwrapReturn(evaluated)
	case *object.BuiltIn:
		return fn.Fn(args...)
	default:
		return &object.Error{Msg: "Not a function object"}

	}

}

func unwrapReturn(obj object.Object) object.Object {
	if ret, ok := obj.(*object.Return); ok {
		return ret.Value
	}
	return obj

}

func extendEnv(fn *object.Function, args []object.Object) *object.Enviorment {
	env := object.SetNewEnclosedEnv(fn.Env)

	for i, p := range fn.Params {
		env.Put(p.Value, args[i])
	}

	return env

}

func evalArguments(args []ast.Expression, env *object.Enviorment) []object.Object {
	var result []object.Object
	for _, e := range args {
		x := Eval(e, env)
		result = append(result, x)
	}

	return result
}

func evalBlockStatements(bstmt *ast.BlockStatement, env *object.Enviorment) object.Object {
	var result object.Object

	for _, stmt := range bstmt.Statements {
		result := Eval(stmt, env)

		if result != nil && result.Type() == object.RETURN_OBJ {
			return result
		}
		if result.Type() == object.ERROR_OBJ && result != nil {
			return result
		}
	}
	return result
}

// func evalProgram(prog *ast.Program, env *object.Enviorment) object.Object {
// 	var result object.Object
// 	for _, stmt := range prog.Statements {
// 		result := Eval(stmt, env)

// 		if ret, ok := result.(*object.Return); ok {
// 			return ret.Value
// 		}

// 	}

// 	return result
// }

func helper(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(stmts []ast.Statement, env *object.Enviorment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)
		if rt, ok := result.(*object.Return); ok {
			return rt.Value
		}

		if result.Type() == object.ERROR_OBJ && result != nil {
			return result
		}
	}

	return result

}
