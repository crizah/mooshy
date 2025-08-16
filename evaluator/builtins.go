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

	"sum": &object.BuiltIn{
		Fn: func(args ...object.Object) object.Object {
			// args all need to be of same type. can be string or int. not bool
			s := len(args)
			if s == 0 {
				return &object.Error{Msg: "No arguments passed"}
			}
			if s == 1 {

				return args[0]

			}

			yeah := true
			t := args[0].Type()

			for _, arg := range args {
				if !yeah {
					return &object.Error{Msg: "Cannot Sum objects of different types"}
				}
				if arg.Type() != t {
					yeah = false
				}
			}

			switch arg := args[0].(type) {
			case *object.Integer:
				s := int64(0)
				for _, a := range args {
					if x, ok := a.(*object.Integer); ok {
						s = s + x.Value
					}

				}

				return &object.Integer{Value: s}
			case *object.String:
				s := ""
				for _, obj := range args {
					if str, ok := obj.(*object.String); ok {
						s = s + str.Value
					}

				}
				return &object.String{Value: s}
			default:

				return &object.Error{Msg: "Cant sum values of type " + arg.Inspect()}

			}

		},
	},
}
