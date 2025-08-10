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

	return Eval(prog) // input ast.Node and return object.Object
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
