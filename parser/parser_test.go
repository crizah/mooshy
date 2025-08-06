package parser

import (
	"mooshy/ast"
	"mooshy/lexer"
	"mooshy/token"
	"testing"
)

// func TestLetStatements(t *testing.T) {
// 	input :=
// 		` let x= 100;
// 	let yay = 8997;
// 	return x;`

// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram() // type arr of Statement
// 	if program == nil {
// 		t.Fatalf("ParseProgram() returned nil")
// 	}
// 	if len(program.Statements) != 3 {
// 		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
// 			len(program.Statements))
// 	}

// 	// tests:= []string{"x", "y", "yay"}

// 	tests := []struct {
// 		expectedIdentifier string
// 	}{
// 		{"x"},
// 		{"yay"},
// 		{"x"},
// 	}

// 	for i, tt := range tests {
// 		// st := program.Statements[i]
// 		// if !testLetStatement(t, st, tt.expectedIdentifier) {
// 		// 	return
// 		// }
// 		st := program.Statements[i]
// 		switch st.TokenLiteral() {
// 		case "let": // already checks thats why its a swicth conditoon
// 			if !testLetStatement(t, st, tt.expectedIdentifier) {
// 				return
// 			}
// 		default:
// 			if !testReturnStatement(t, st, tt.expectedIdentifier) {
// 				return
// 			}

// 		}

// 	}
// }

func testLetStatement(t *testing.T, st ast.Statement, exName string, expected interface{}) bool {
	if st.TokenLiteral() != "let" {
		t.Errorf("st.TokenLiteral not 'let'. got=%q", st.TokenLiteral())
		return false
	}

	letSt, ok := st.(*ast.LetStatement) // bool check if it is a let statement but abive already does that with its liyeral but can be a bug
	if !ok {
		t.Errorf("st not *ast.LetStatement. got=%T", st)
		return false
	}

	if letSt.Name.Value != exName { // needs to be expected expression
		t.Errorf("letSt.Name.Value not '%s'. got=%s", exName, letSt.Name.Value)
		return false
	}

	if letSt.Name.TokenLiteral() != exName {
		t.Errorf("letSt.Name not '%s'. got=%s", exName, letSt.Name)
		return false
	}

	if !testLiteralExpression(t, letSt.Value, expected) {
		return false
	}

	return true

}

func testReturnStatement(t *testing.T, st ast.Statement, exName string) bool {
	if st.TokenLiteral() != "return" {
		t.Errorf("st.TokenLiteral not 'return'. got=%q", st.TokenLiteral())
		return false
	}

	retStmt, ok := st.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("not of type Return statement. got %T", st)
		return false
	}

	if retStmt.TokenLiteral() != "return" {
		t.Errorf("retStmt.TokenLiteral not 'return'. got=%q", retStmt.TokenLiteral())
		return false
	}
	if retStmt.Value.Value != exName {
		t.Errorf("retStmt.Value.Value not '%s'. got=%s", exName, retStmt.Value.Value)
		return false
	}
	return true

}

func TestExpression(t *testing.T) {
	input := `5;`
	l := lexer.New(input)
	// fmt.Printf("check1")
	p := New(l)
	// fmt.Printf("check2\n")

	program := p.ParseProgram()
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 2 statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("program.Statements is not ExpressionStatement. got %T ", program.Statements[0])
		return
	}
	// fmt.Printf("check3\n")
	ident, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("stmt.Expression is not IntegerLiteral. got %T", stmt.Expression)
		return
	}

	if ident.Token.Type != token.INT {
		t.Errorf("ident.Token.Type not INT. got=%q", ident.Token.Type)
		return
	}

	if ident.Value != 5 {
		t.Error("ident.Value not 5. got=", ident.Value)
		return
	}

	if ident.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not '5'. got=%q", ident.TokenLiteral())
		return
	}

}

func TestPrefix(t *testing.T) {
	// input := `
	// !5;
	// -10;
	// `
	tests := []struct {
		input    string
		operator string
		value    int64
	}{
		{"!5", "!", 5},
		{"-10", "-", 10},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements has wrong len. input=%q", tt.input)
			return
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("stmt not ExpressionStatement. got %T", stmt)
			return
		}

		st, ok := stmt.Expression.(*ast.PrefixExpression)
		// PrefixExpression needs to have field Operator of type string and Right of type Expression wh ich is of integerLiteral, Expression
		if !ok {
			t.Errorf("stmt not PrefixExpression. got %T", stmt)
			return
		}

		if st.Operator != tt.operator {
			t.Errorf("stmt.Operator not %s. got %s", tt.operator, st.Operator)
			return
		}

		if !testIntegerLiteral(t, st.Right, tt.value) {
			return
		}

	}
}

func TestInfix(t *testing.T) {
	tests := []struct {
		input    string
		left     int64
		right    int64
		operator string
	}{ // have to do this with strings too, str1+str2 so Idetifier
		{"5 +10", 5, 10, "+"},
		{"7- 9", 7, 9, "-"},
		{"2* 8", 2, 8, "*"},
		{"10/ 5", 10, 5, "/"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()

		if len(program.Statements) != 1 {
			t.Errorf("program.Statements has wrong len. input=%q", tt.input)
			return
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("program.Statements[0] not ExpressionStatement. got %T", program.Statements)
			return
		}

		infx, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Errorf("stmt not InfixExpression. got %T", stmt)
			return
		}

		if infx.Operator != tt.operator {
			t.Errorf("stmt.Operator not %s. got %s", tt.operator, infx.Operator)
			return
		}

		if !testIntegerLiteral(t, infx.Left, tt.left) {
			return
		}
		if !testIntegerLiteral(t, infx.Right, tt.right) {
			return
		}

	}
}
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	intLit, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("not IntegerLiteral. got %T", il)
		return false
	}

	if intLit.Value != value {
		t.Errorf("intLit.Value not as expectedValue. got %d, want %d", intLit.Value, value)
		return false
	}

	// change int64 to string

	return true

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a)*b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a+b)+c)",
		},
		{
			"a + b - c",
			"((a+b)-c)",
		},
		{
			"a * b * c",
			"((a*b)*c)",
		},
		{
			"a * b / c",
			"((a*b)/c)",
		},
		{
			"a + b / c",
			"(a+(b/c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a+(b*c))+(d/e))-f)",
		},
		// {
		// 	"3 + 4; -5 * 5",
		// 	"(3+4)((-5)*5)",
		// },
		{
			"5 > 4 == 3 < 4",
			"((5>4)==(3<4))",
		},
		{
			"5 < 4 != false",
			"((5<4)!=false)",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3+(4*5))==((3*1)+(4*5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3+(4*5))==((3*1)+(4*5)))",
		},
		{"- 1 + 2",
			"((-1)+2)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		// checkParserErrors(t, p)
		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, expected string) bool {
	id, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("not Identifier. got %T", exp)
		return false
	}

	if id.Value != expected {
		t.Errorf("id.Value not as expectedValue. got %s, want %s", id.Value, expected)
		return false

	}
	if id.TokenLiteral() != expected {
		t.Errorf("id.TokenLiteral not as expectedValue. got %s, want %s", id.TokenLiteral(), expected)
		return false
	}

	return true

}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolExpression(t, exp, v)

	}
	t.Errorf("type of interface not defined")

	return false

}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	infx, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("not InfixExpression. got %T", exp)
		return false
	}

	if !testLiteralExpression(t, infx.Left, left) {
		t.Errorf("infx.left not as expected. got %T, want %T", infx.Left, left)
		return false
	}

	if !testLiteralExpression(t, infx.Right, right) {
		t.Errorf("infx.right not as expected. got %T, want %T", infx.Right, right)
		return false
	}

	if infx.Operator != operator {
		t.Errorf("infx.Operator not as expected. got %s, want %s", infx.Operator, operator)
		return false
	}

	return true

}

func TestYeah(t *testing.T) {

	input := " meow + purr"

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("not ExpressionStatement. got %T", prog.Statements[0])
		return
	}

	if !testInfixExpression(t, stmt.Expression, "meow", "+", "purr") {
		t.Errorf("Test failed")
		return
	}
}

func testBoolExpression(t *testing.T, exp ast.Expression, expectedBool bool) bool {
	boolExp, ok := exp.(*ast.BoolExpression)
	if !ok {
		t.Errorf("not BoolExpression. got %T", exp)
		return false
	}

	if boolExp.Value != expectedBool {
		t.Errorf("boolExp.Value not as expected. got %v, want %v", boolExp.Value, expectedBool)
		return false
	}

	return true

}

func TestNaur(t *testing.T) {
	input := `
	true;
	false;
	let x = true;`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	if len(prog.Statements) != 3 {
		t.Errorf("expected 3 statements, got %d", len(prog.Statements))
		return
	}

	stmt1, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("not ExpressionStatement. got %T", prog.Statements[0])
		return
	}

	stmt2, ok2 := prog.Statements[1].(*ast.ExpressionStatement)

	if !ok2 {
		t.Errorf("not ExpressionStatement. got %T", prog.Statements[1])
		return
	}

	stmt3, ok3 := prog.Statements[2].(*ast.LetStatement)

	if !ok3 {
		t.Errorf("not LetStatement. got %T", prog.Statements[2])
		return
	}

	if !ok {
		t.Errorf("not ExpressionStatement. got %T", prog.Statements[0])
		return
	}

	if !testBoolExpression(t, stmt1.Expression, true) {
		t.Errorf("bool expression Test failed")
		return
	}

	if !testBoolExpression(t, stmt2.Expression, false) {
		t.Errorf("bool expression Test failed")
		return
	}

	if !testLetStatement(t, stmt3, "x", true) {
		t.Errorf("let statement Test failed")
		return
	}

}

func TestHmm(t *testing.T) {
	input := ` 4<5 != false;`

	tests := []struct {
		input    string
		expected bool
		exp      string
	}{
		{input, false, "((4<5)!=false)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		prog := p.ParseProgram()

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("not expression statement. got %T", stmt)
			return
		}

		infx, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Errorf("not infix expression. got %T", infx)
			return
		}
		if !testInfixExpression(t, infx.Left, 4, "<", 5) {
			t.Errorf("test not passed here")
		}
		// if !testInfixExpression(t, infx, infx.Left, "!=", false) { // not wotrking here
		// 	t.Errorf("test not passed HERE")
		// 	return

		// }

		if !testBoolExpression(t, infx.Right, false) {
			t.Errorf("not passed")
			return
		}

	}
}

func TestThistoo(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"(5 + 5) * 2",
			"((5+5)*2)",
		},
		{
			"-(5 + 5)",
			"(-(5+5))",
		},
		{
			"!(true == true)",
			"(!(true==true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}

	}
}

func TestIfStatements(t *testing.T) {
	input := `

	if (5<4){
	    let x = 10;
	    let y = 20;

	}else{
	    let z = 10;
	    let g = 20;
	}`

	l := lexer.New(input)
	p := New(l)

	prog := p.ParseProgram()
	if len(prog.Statements) != 1 {
		t.Errorf("len of stetemets not as expected. expected %d, got %d", 1, len(prog.Statements))
		return
	}
	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("statement not an expression statement. got %T", stmt)
		return
	}

	ifExp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Errorf("not an ifExpression. got %T", ifExp)
		return
	}

	cond, ok := ifExp.Condition.(*ast.InfixExpression) // for this conditon is (5<4)
	if !ok {
		t.Errorf("condition is not InfixExpression. got %T", cond)
		return
	}

	if !testInfixExpression(t, cond, 5, "<", 4) {
		t.Errorf("got error")
		return
	}

	// cons, ok := ifExp.Consequence.(*ast.BlockStatement)
	// if !ok{
	// 	t.Errorf("concequence is not a BlockStatement. got %T", cons)
	// 	return
	// }

	if len(ifExp.Consequence.Statements) != 2 {
		t.Errorf("len of consequence not as expected. expected %d, got %d", 2, len(ifExp.Consequence.Statements))
		return
	}

	if !testLetStatement(t, ifExp.Consequence.Statements[0], "x", 10) {
		t.Errorf("not passed")
		return
	}

	if !testLetStatement(t, ifExp.Consequence.Statements[1], "y", 20) {
		t.Errorf("not passed")
		return
	}

	if ifExp.Alternative == nil {
		t.Errorf("alternative not available")
		return
	}

	if len(ifExp.Alternative.Statements) != 2 {
		t.Errorf("len of consequence not as expected. expected %d, got %d", 2, len(ifExp.Alternative.Statements))
		return
	}

	if !testLetStatement(t, ifExp.Alternative.Statements[0], "z", 10) {
		t.Errorf("not passed")
		return
	}

	if !testLetStatement(t, ifExp.Alternative.Statements[1], "g", 20) {
		t.Errorf("not passed")
		return
	}

}

func TestFunctions(t *testing.T) {
	input := `
	func(x, y){
	let z = 20;
	return x; }`

	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()
	// if prog.String() != input {
	// 	t.Errorf("not as expected. expected %s. got %s", input, prog.String())
	// 	return

	// }
	// expression statement and function literal

	if len(prog.Statements) != 1 {
		t.Errorf("len of statements worng. expected %d, got %d", 1, len(prog.Statements))
		return
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("not an expression statement. got %T", stmt)
		return
	}

	fexp, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Errorf("not a function statement. got %T", fexp)
		return
	}

	if len(fexp.Parameter) != 2 {
		t.Errorf("len of parameters is not 2. got %d", len(fexp.Parameter))
		return

	}
	if !testIdentifier(t, fexp.Parameter[0], "x") {
		t.Errorf("naur")
		return
	}

	if !testIdentifier(t, fexp.Parameter[1], "y") {
		t.Errorf("did not pass")
		return
	}

	if len(fexp.Body.Statements) != 2 {
		t.Errorf("number of ststemnets in block not 2. got %d", len(fexp.Body.Statements))
		return
	}

	if !testLetStatement(t, fexp.Body.Statements[0], "z", 20) {
		t.Errorf("did not pass")
		return
	}

	if !testReturnStatement(t, fexp.Body.Statements[1], "x") {
		t.Errorf("didnt pass")
		return
	}

}

// func TestFunctionLiterals(t *testing.T){
// 	input:= `
// 	let add = func(x, y){
// 	let z = 20;
// 	return x+y+20; }`

// 	l:= lexer.New(input)
// 	p:= New(l)
// 	prog := p.ParseProgram()

// 	if len(prog.Statements)!=1{
// 		t.Errorf("len of statements worng. expected %d, got %d", 1, len(prog.Statements))
// 		return
// 	}

// 	stmt, ok := prog.Statements[0].(*ast.LetStatement)
// 	if !ok {
// 		t.Errorf("stmt not let statement. got %T", stmt)
// 		return
// 	}

// 	// if !testLetStatement(t, stmt, "add", float32(3.5)){
// 	// 	t.Errorf("test not passed")
// 	// 	return
// 	// }

// }
