package evaluator

import (
	"mooshy/ast"
	"mooshy/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Bool{Value: true}
	FALSE = &object.Bool{Value: false}
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
		return NULL
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
		return NULL

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
		return NULL
	}

}

func evalInfixOp(left object.Object, right object.Object, operator string) object.Object {
	// REFINE THIS WTF EVEN IS THIS CODE.

	// both left and right dont NEED rto be integres. can be of string type and return will be concacted string.
	// but for now, keep them both as only integers allowed

	if right.Type() == object.INTEGER_OBJ && left.Type() == object.INTEGER_OBJ {
		// also do if of differetent object types and == and != operators instead of null

		r, ok := right.(*object.Integer)
		if !ok {
			return NULL
		}

		l, ok := left.(*object.Integer)
		if !ok {
			return NULL
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
			return NULL
		}
	} else if right.Type() == object.BOOL_OBJ && left.Type() == object.BOOL_OBJ {

		r, ok := right.(*object.Bool)
		if !ok {
			return NULL
		}

		l, ok := left.(*object.Bool)
		if !ok {
			return NULL
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
		}

	}

	return NULL

}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.BoolExpression:
		return helper(node.Value) // so that we dont have to make a new instance everytime
		// return &object.Bool{Value: node.Value}

	// case *ast.StringLiteral:
	// 	return &object.String{Value: node.Value}
	case *ast.PrefixExpression: // can be - or !
		right := Eval(node.Right)
		return evaluateOp(right, node.Operator)
	case *ast.InfixExpression:
		right := Eval(node.Right)
		left := Eval(node.Left)
		return evalInfix(right, left, node.Operator)
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	}

	return nil
}

func helper(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result

}
