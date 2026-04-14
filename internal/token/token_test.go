package token

import (
	"testing"
)

func TestLookupIdent(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"fn", FUNCTION},
		{"let", LET},
		{"true", TRUE},
		{"false", FALSE},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"while", WHILE},
		{"null", NULL},
		// 非关键字测试
		{"x", IDENT},
		{"foobar", IDENT},
		{"function", IDENT}, // "function" 不是关键字，我们的关键字是 "fn"
	}

	for _, tt := range tests {
		tok := LookupIdent(tt.input)
		if tok != tt.expected {
			t.Errorf("LookupIdent(%q) wrong. expected=%q, got=%q", tt.input, tt.expected, tok)
		}
	}
}
