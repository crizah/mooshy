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

	// if !testLiteralExpression(t, retStmt.Value, exName) {
	// 	t.Errorf("did not pass")
	// 	return false
	// }

	if retStmt.Value.String() != exName {
		t.Errorf("retStmt.Value.Value not '%s'. got=%s", exName, retStmt.Value.String())
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
		t.Fatalf("program.Statements does not contain 1 statements. got=%d",
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

		{
			"a + add(b * c) + d",
			"((a+add((b*c)))+d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2*3), (4+5), add(6, (7*8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a+b)+((c*d)/f))+g))",
		},
		{
			"x++",
			"(x++)",
		},
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
	let z = "meow";
	return x+y; }`

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
		// t.Errorf("statements are %T, %T, %T, %T", fexp.Body.Statements[0], fexp.Body.Statements[1], fexp.Body.Statements[2], fexp.Body.Statements[3])
		return
	}

	// testing the let statment

	letStmt, ok := fexp.Body.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Errorf("not a let statement. got %T", letStmt)
		return
	}

	if letStmt.Name.String() != "z" {
		t.Errorf("name of let statement is not z. got %s", letStmt.Name.String())
		return
	}

	str, ok := letStmt.Value.(*ast.StringLiteral)
	if !ok {
		t.Errorf("not a string literal. got %T", str)
		return
	}

	if str.Value != "meow" {
		t.Errorf("value of string literal is not meow. got %s", str.Value)
		return
	}

	// if !testLetStatement(t, fexp.Body.Statements[0], "z", 20) {
	// 	t.Errorf("did not pass")
	// 	return
	// }

	rtstmt, ok := fexp.Body.Statements[1].(*ast.ReturnStatement)
	if !ok {
		t.Errorf("not a return statement. got %T", rtstmt)
		return
	}

	if !testReturnStatement(t, rtstmt, "(x+y)") {
		t.Errorf("didnt pass")
		return
	}

}

func TestCallExpression(t *testing.T) {
	input := `
	add(x+y, c+d, 7);`
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	if len(prog.Statements) != 1 {
		t.Errorf("len not 1. got %d", len(prog.Statements))

	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("not expresion statement. got %T", stmt)

	}

	callEXp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Errorf("not call expression. got %T", stmt.Expression)

	}

	if !testIdentifier(t, callEXp.Function, "add") {
		t.Errorf("didnt pass")
		return

	}

	if len(callEXp.Arguments) != 3 {
		t.Errorf("len not 3. got %d", len(callEXp.Arguments))

	}

	a, ok := callEXp.Arguments[0].(*ast.InfixExpression)
	if !ok {
		t.Errorf("not an infix expression. got %T", a)
		return
	}

	b, ok := callEXp.Arguments[1].(*ast.InfixExpression)
	if !ok {
		t.Errorf("not an infix expression. got %T", b)
		return
	}

	if !testInfixExpression(t, a, "x", "+", "y") {
		t.Errorf("didnt pass")
		return
	}

	if !testInfixExpression(t, b, "c", "+", "d") {
		t.Errorf("didnt pass")
		return
	}

	if !testIntegerLiteral(t, callEXp.Arguments[2], int64(7)) {
		t.Errorf("didnt pass")
		return
	}
}

func TestArrays(t *testing.T) {
	input := `
     [5, 6, 987];`
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	if len(prog.Statements) != 1 {
		t.Errorf("len not 1. got %d", len(prog.Statements))

	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("not expresion statement. got %T", stmt)

	}

	arrExp, ok := stmt.Expression.(*ast.ArrayExpression)
	if !ok {
		t.Errorf("not array expression. got %T", stmt.Expression)
		return

	}

	if len(arrExp.Value) != 3 {
		t.Errorf("len not 3. got %d", len(arrExp.Value))

	}

	a, ok := arrExp.Value[0].(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("not an IntegerLiteral. got %T", a)
		return
	}

	if !testIntegerLiteral(t, a, int64(5)) {
		t.Errorf("didnt pass")
		return
	}

}

func TestIndexExp(t *testing.T) {
	input := `
     arr[5];`
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	if len(prog.Statements) != 1 {
		t.Errorf("len not 1. got %d", len(prog.Statements))

	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("not expresion statement. got %T", stmt)

	}

	arrExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Errorf("not index expression. got %T", stmt.Expression)
		return

	}

	if !testIdentifier(t, arrExp.Name, "arr") {
		t.Errorf("got error")
		return
	}

	if arrExp.Index != int64(5) {
		t.Errorf("not 5 got %d", arrExp.Index)
		return
	}

}

func TestReassignExp(t *testing.T) {
	input := `
    x = 12;`
	l := lexer.New(input)
	p := New(l)
	prog := p.ParseProgram()

	if len(prog.Statements) != 1 {
		t.Errorf("len not 1. got %d", len(prog.Statements))

	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Errorf("not expresion statement. got %T", stmt)

	}

	arrExp, ok := stmt.Expression.(*ast.ReAssignExpression)
	if !ok {
		t.Errorf("not Reassign expression. got %T", stmt.Expression)
		return

	}

	if !testIdentifier(t, arrExp.Name, "x") {
		t.Errorf("got error")
		return
	}

	if !testIntegerLiteral(t, arrExp.Value, int64(12)) {
		t.Errorf("got error")
		return

	}

}

func testPostFix(t *testing.T, exp ast.Expression, left interface{}, operator string) bool {
	infx, ok := exp.(*ast.PostfixExpression)
	if !ok {
		t.Errorf("not PostfixExpression. got %T", exp)
		return false
	}

	if !testLiteralExpression(t, infx.Left, left) {
		t.Errorf("post.left not as expected. got %T, want %T", infx.Left, left)
		return false
	}

	if infx.Operator != operator {
		t.Errorf("post.Operator not as expected. got %s, want %s", infx.Operator, operator)
		return false
	}

	return true

}

func TestPostFixExp(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		Literal  interface{}
	}{
		{"12++", "++", 12},
		{"a++", "++", "a"},
		{"a--", "--", "a"},
		{"12--", "--", 12},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		prog := p.ParseProgram()
		if len(prog.Statements) != 1 {
			t.Errorf("length not 1. got %d", len(prog.Statements))
			return
		}

		exp, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("not exp statement got %T", exp)

		}

		post, ok := exp.Expression.(*ast.PostfixExpression)

		if !ok {
			t.Errorf("stmt not PrefixExpression. got %T", post)
			return
		}

		if !testPostFix(t, post, tt.Literal, tt.operator) {
			t.Errorf("didnt pass")
			return
		}
	}

}

func TestForLoops(t *testing.T) {
	tests := []struct {
		input string
	}{
		{`for(let i =0; i<12; i++){ 
		let x = 12;
		x++;
}`},
		// start = let i =0(let stmt), i =0 (assign expression stmt), i (ident expression statement )
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		prog := p.ParseProgram()

		if len(prog.Statements) != 1 {
			t.Errorf("length not 1. got %d", len(prog.Statements))
			return
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("not expression statement. got %T", stmt)
			return
		}

		exp, ok := stmt.Expression.(*ast.ForLoopExpressions)
		if !ok {
			t.Errorf("not for loop expression. got %T", exp)
			return
		}

		para, ok := exp.Params.(*ast.BlockExpression)
		if !ok {
			t.Errorf("not block exp. got %T", para)
			return

		}

		testLetStatement(t, para.Start, "i", 0)
		testInfixExpression(t, para.Condition, "i", "<", 12)
		testPostFix(t, para.Iterator, "i", "++")

		if len(exp.Body.Statements) != 2 {
			t.Errorf("Not 2. got %d", len(exp.Body.Statements))
			// t.Errorf("%T, %T, %T", exp.Body.Statements[0], exp.Body.Statements[1], exp.Body.Statements[2])
			return
		}

		testLetStatement(t, exp.Body.Statements[0], "x", 12)

		yeah, ok := exp.Body.Statements[1].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("not exp statement got %T", exp)

		}

		post, ok := yeah.Expression.(*ast.PostfixExpression)
		if !ok {
			t.Errorf("Not postFix exp. got %T", post)
		}

		testPostFix(t, post, "x", "++")

		// testLetStatement(t, )

		// LET STATEMENT CANT BE EXPReSSION

	}
}
