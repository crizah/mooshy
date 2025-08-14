package evaluator

import (
	"mooshy/lexer"
	"mooshy/object"
	"mooshy/parser"
	"testing"
)

func TestIntegerObjects(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"876", 876},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestBooleanObjects(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"false", false},
		{"true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

// func TestStringObjects(t *testing.T) {
// 	tests := []struct {
// 		input    string
// 		expected string
// 	}{
// 		{"meow", `"meow"`},
// 		{"purr", `"purr"`},
// 	}

// 	for _, tt := range tests {
// 		evaluated := testEval(t, tt.input)
// 		testStringObject(t, evaluated, tt.expected)
// 	}

// }

// func testStringObject(t *testing.T, obj object.Object, expected string) bool {
// 	str, ok := obj.(*object.String)
// 	if !ok {
// 		t.Errorf("object is not String. got=%T", obj)
// 		return false
// 	}

// 	if str.Value != expected {
// 		t.Errorf("value not as expected. expected %s, got %s", expected, str.Value)
// 		return false
// 	}

// 	return true

// }
func testEval(t *testing.T, input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)

	prog := p.ParseProgram()
	env := object.NewEnv()

	return Eval(prog, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expectedValue int64) bool {
	it, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T", obj)
		return false
	}

	if it.Value != expectedValue {
		t.Errorf("value not as expected. expected %d, got %d", expectedValue, it.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expectedValue bool) bool {
	bl, ok := obj.(*object.Bool)
	if !ok {
		t.Errorf("not boolean object. got=%T", obj)
		return false
	}

	if bl.Value != expectedValue {
		t.Errorf("value not as expected. expected %t, got %t", expectedValue, bl.Value)
		return false
	}

	return true
}

func TestPrefixNot(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!!true", true},
		{"!!false", false},
		{"!5", false}, // 5 is considered truthly
		{"!!5", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestPrefixMinus(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"-10", -10},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)

	}

}

func TestInfix(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestFuncObjects(t *testing.T) {
	// tests := []struct {
	// 	input    string
	// 	expected interface{}
	// }{
	// 	{"if (true) { 10 }", 10},
	// 	{"if (false) { 10 }", nil},
	// 	{"if (1) { 10 }", 10},
	// 	{"if (1 < 2) { 10 }", 10},
	// 	{"if (1 > 2) { 10 }", nil},
	// 	{"if (1 > 2) { 10 } else { 20 }", 20},
	// 	{"if (1 < 2) { 10 } else { 20 }", 10},
	// }

	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { return 10; }", 10},
		{"if (false) { return 10; }", nil},
		{"if (1) { return 10; }", 10},
		{"if (1 < 2) { return 10; }", 10},
		{"if (1 > 2) {return 10;}", nil},
		{"if (1 > 2) { return 10; } else { return 20; }", 20},
		{"if (1 < 2) { return 10; } else { return 20;}", 10}, // needs to be followed by a return statement
		{`if(10 > 1){ if(9 > 1) {
		return 10; } return 1; }`, 10},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		i, ok := tt.expected.(int)

		if ok {
			testIntegerObject(t, evaluated, int64(i))
		} else {
			testNullObject(t, evaluated)
		}

	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("not a null object got %T", obj)
		return false

	}

	return true

}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10}, // only return till the first ;
		{"9; return 2 * 5; 9;", 10},
	}
	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestReturn(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{{"return 98;", 98},
		{"return false;", false},
		{`return "meow";`, "meow"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		i, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(i))
		}

		b, ok := tt.expected.(bool)
		if ok {
			testBooleanObject(t, evaluated, b)
		}

		s, ok := tt.expected.(string)
		if ok {
			testStringObjects(t, evaluated, s)
		}

	}
}

func TestStringObject(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"meow"`, "meow"},
	}

	for _, tt := range tests {
		evalauted := testEval(t, tt.input)
		testStringObjects(t, evalauted, tt.expected)
	}
}

func testStringObjects(t *testing.T, obj object.Object, expected string) bool {
	str, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T", obj)
		return false
	}

	if str.Value != expected {
		t.Errorf("not as expected. expected %s. got %s", expected, str.Value)
		return false
	}

	return true
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"let x = 10; x;", 10},
		{"let a = false; a;", false},
		{"let b = 2 * 5; b;", 10},
		{"let x = 10>1; x;", true},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		{"let a = 5; let b = a; b;", 5},
		{`let a = "meow"; let b = "wwww"; let c = a + b + "!!"; c;`, "meowwwww!!"},
	}

	for _, tt := range tests {
		evaluated := testEval(t, tt.input)
		val, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(val))
		} else {
			b, ok := tt.expected.(bool)
			if ok {
				testBooleanObject(t, evaluated, b)
			} else {
				testStringObjects(t, evaluated, tt.expected.(string))
			}
		}

	}
}

func TestFunctions(t *testing.T) {
	input := "func(x) { (x+2); }"
	evaluated := testEval(t, input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Errorf("Not a function object. got %T", fn)
		return

	}

	if len(fn.Params) != 1 {
		t.Errorf("length of params not 1. got %d", len(fn.Params))
		return

	}

	if fn.Params[0].Value != "x" {
		t.Errorf("not as expected. expected %s. got %s", "x", fn.Params[0].Value)
		return
	}

	if len(fn.Body.Statements) != 1 {
		t.Errorf("length of body not 1. got %d", len(fn.Body.Statements))
		return
	}

	if fn.Body.Statements[0].String() != "(x+2)" {
		t.Errorf("not as expected. got %s", fn.Body.Statements[0].String())
		return
	}

}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		// HAS TO HAVE A RETURN STATEMENT . NEEDS TO RETURN A VALUE

		{"let identity = func(x) { return x; }; identity(5);", 5},
		{"let double = func(x) { return (x * 2); }; double(5);", 10},
		{"let add = func(x, y) { return (x + y); }; add(5, 5);", 10},
		{"let add = func(x, y) { return (x + y); }; add(5 + 5, add(5, 5));", 20},
		{"func(x) { return x; }(5)", 5}, // ???
	}
	for _, tt := range tests {
		testIntegerObject(t, testEval(t, tt.input), tt.expected)
	}
}
