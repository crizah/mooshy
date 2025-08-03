package parser

import (
	"mooshy/ast"
	"mooshy/lexer"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input :=
		` let x= 100;
	let yay = 8997;
	return x;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram() // type arr of Statement
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	// tests:= []string{"x", "y", "yay"}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"yay"},
		{"x"},
	}

	for i, tt := range tests {
		// st := program.Statements[i]
		// if !testLetStatement(t, st, tt.expectedIdentifier) {
		// 	return
		// }
		st := program.Statements[i]
		switch st.TokenLiteral() {
		case "let": // already checks thats why its a swicth conditoon
			if !testLetStatement(t, st, tt.expectedIdentifier) {
				return
			}
		default:
			if !testReturnStatement(t, st, tt.expectedIdentifier) {
				return
			}

		}

	}
}

func testLetStatement(t *testing.T, st ast.Statement, exName string) bool {
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
