package parser

import (
	"fmt"
	"testing"

	"make_a_lang/internal/ast"
	"make_a_lang/internal/lexer"
	"make_a_lang/internal/token"
)

func TestParseLeftAssignStmts(t *testing.T) {
	l := lexer.NewLexer("let x = 5; return 10;")
	p := NewParser(l)
	root := p.ParseProgram()

	expectStmts := []ast.Statement{
		&ast.LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let", Line: 1, Column: 1, EndLine: 1, EndColumn: 4, StartOffset: 0, EndOffset: 3},
			Name: &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "x", Line: 1, Column: 5, EndLine: 1, EndColumn: 6, StartOffset: 4, EndOffset: 5},
				Value: "x",
			},
			Value: &ast.IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "5", Line: 1, Column: 9, EndLine: 1, EndColumn: 10, StartOffset: 8, EndOffset: 9},
				Value: 5,
			},
		},
		&ast.ReturnStatement{
			Token: token.Token{Type: token.RETURN, Literal: "return", Line: 1, Column: 12, EndLine: 1, EndColumn: 18, StartOffset: 11, EndOffset: 17},
			ReturnValue: &ast.IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "10", Line: 1, Column: 19, EndLine: 1, EndColumn: 21, StartOffset: 18, EndOffset: 20},
				Value: 10,
			},
		},
	}

	if len(root.Statements) != len(expectStmts) {
		t.Fatalf("len(root.Statements) = %d, want %d", len(root.Statements), len(expectStmts))
	}

	for i, stmt := range root.Statements {
		if stmt.TokenLiteral() != expectStmts[i].TokenLiteral() {
			t.Errorf("stmt.TokenLiteral() = %q, want %q", stmt.TokenLiteral(), expectStmts[i].TokenLiteral())
		}
	}
}

func TestParseInfixExpressions(t *testing.T) {
	l := lexer.NewLexer("5 + 5; 10 - 2;")
	p := NewParser(l)
	root := p.ParseProgram()

	expectStmts := []ast.Statement{
		&ast.ExpressionStatement{
			Token: token.Token{Type: token.INT, Literal: "5", Line: 1, Column: 1, EndLine: 1, EndColumn: 2, StartOffset: 0, EndOffset: 1},
			Expression: &ast.InfixExpression{
				Token:    token.Token{Type: token.PLUS, Literal: "+", Line: 1, Column: 3, EndLine: 1, EndColumn: 4, StartOffset: 2, EndOffset: 3},
				Left:     &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5", Line: 1, Column: 1, EndLine: 1, EndColumn: 2, StartOffset: 0, EndOffset: 1}, Value: 5},
				Operator: "+",
				Right:    &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5", Line: 1, Column: 5, EndLine: 1, EndColumn: 6, StartOffset: 4, EndOffset: 5}, Value: 5},
			},
		},
		&ast.ExpressionStatement{
			Token: token.Token{Type: token.INT, Literal: "10", Line: 1, Column: 8, EndLine: 1, EndColumn: 10, StartOffset: 7, EndOffset: 9},
			Expression: &ast.InfixExpression{
				Token:    token.Token{Type: token.MINUS, Literal: "-", Line: 1, Column: 11, EndLine: 1, EndColumn: 12, StartOffset: 10, EndOffset: 11},
				Left:     &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "10", Line: 1, Column: 8, EndLine: 1, EndColumn: 10, StartOffset: 7, EndOffset: 9}, Value: 10},
				Operator: "-",
				Right:    &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "2", Line: 1, Column: 13, EndLine: 1, EndColumn: 14, StartOffset: 12, EndOffset: 13}, Value: 2},
			},
		},
	}

	if len(root.Statements) != len(expectStmts) {
		t.Fatalf("len(root.Statements) = %d, want %d", len(root.Statements), len(expectStmts))
	}

	for i, stmt := range root.Statements {
		if stmt.TokenLiteral() != expectStmts[i].TokenLiteral() {
			t.Errorf("stmt.TokenLiteral() = %q, want %q", stmt.TokenLiteral(), expectStmts[i].TokenLiteral())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) bool {
	il, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp not *ast.IntegerLiteral. got=%T", exp)
		return false
	}

	if il.Value != value {
		t.Errorf("il.Value = %d, want %d", il.Value, value)
		return false
	}

	if il.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("il.TokenLiteral() = %q, want %q", il.TokenLiteral(), fmt.Sprintf("%d", value))
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value = %q, want %q", ident.Value, value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() = %q, want %q", ident.TokenLiteral(), value)
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bl, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}

	if bl.Value != value {
		t.Errorf("bl.Value = %t, want %t", bl.Value, value)
		return false
	}

	if bl.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bl.TokenLiteral() = %q, want %q", bl.TokenLiteral(), fmt.Sprintf("%t", value))
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

// 伪代码提示：
func TestLetStatements(t *testing.T) {
    tests := []struct {
        input              string
        expectedIdentifier string
        expectedValue      interface{} // 可以是 int, string, bool 等
    }{
        {"let x = 5;", "x", 5},
		{"let y = true;", "y", true},
		{"let z = false;", "z", false},
		{"let foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		p := NewParser(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.LetStatement. got=%T", program.Statements[0])
		}

		if stmt.Name.Value != tt.expectedIdentifier {
			t.Errorf("stmt.Name.Value = %q, want %q", stmt.Name.Value, tt.expectedIdentifier)
		}

		if !testLiteralExpression(t, stmt.Value, tt.expectedValue) {
			return
		}
	}
}