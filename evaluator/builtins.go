package evaluator

import "mooshy/object"

var builtins = map[string]*object.BuiltIn{
	"len": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return &object.Error{Msg: "Expected 1 arguments"}
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return &object.Error{Msg: "Expected String Object"}

			}

		},
	},
}
