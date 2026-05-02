package lexer_test

import (
	"make_a_lang/internal/lexer"
	token "make_a_lang/internal/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    token.TokenType
		wantLiteral string
	}{
		{
			name:        "simple assignment",
			input:       "let x = 5;",
			wantType:    token.LET,
			wantLiteral: "let",
		},
		{
			name:        "function definition",
			input:       "fun main() {}",
			wantType:    token.FUNCTION,
			wantLiteral: "fun",
		},
		{
			name:        "if statement",
			input:       "if (x < 10) { return true; } else { return false; }",
			wantType:    token.IF,
			wantLiteral: "if",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			got := lexer.NextToken(l)
			if got.Type != tt.wantType || got.Literal != tt.wantLiteral {
				t.Errorf("NextToken() = (%v, %q), want (%v, %q)", got.Type, got.Literal, tt.wantType, tt.wantLiteral)
			}
		})
	}
}

func TestNextTokenPositionConsistency(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  token.Token
	}{
		{
			name:  "single char",
			input: "+",
			want:  token.Token{Type: token.PLUS, Literal: "+", Line: 1, Column: 1, EndLine: 1, EndColumn: 2, StartOffset: 0, EndOffset: 1},
		},
		{
			name:  "double char",
			input: "==",
			want:  token.Token{Type: token.EQ, Literal: "==", Line: 1, Column: 1, EndLine: 1, EndColumn: 3, StartOffset: 0, EndOffset: 2},
		},
		{
			name:  "single assign",
			input: "=",
			want:  token.Token{Type: token.ASSIGN, Literal: "=", Line: 1, Column: 1, EndLine: 1, EndColumn: 2, StartOffset: 0, EndOffset: 1},
		},
		{
			name:  "simple identifier",
			input: "foo",
			want:  token.Token{Type: token.IDENT, Literal: "foo", Line: 1, Column: 1, EndLine: 1, EndColumn: 4, StartOffset: 0, EndOffset: 3},
		},
		{
			name:  "cross line",
			input: "\nfoo",
			want:  token.Token{Type: token.IDENT, Literal: "foo", Line: 2, Column: 1, EndLine: 2, EndColumn: 4, StartOffset: 1, EndOffset: 4},
		},
		{
			name:  "string",
			input: "\"a\nb\"",
			want:  token.Token{Type: token.STRING, Literal: "a\nb", Line: 1, Column: 1, EndLine: 2, EndColumn: 3, StartOffset: 0, EndOffset: 5},
		},
		{
			name:  "comment",
			input: "// hi\n+",
			want:  token.Token{Type: token.PLUS, Literal: "+", Line: 2, Column: 1, EndLine: 2, EndColumn: 2, StartOffset: 6, EndOffset: 7},
		},
		{
			name:  "illegal char",
			input: "@",
			want:  token.Token{Type: token.ILLEGAL, Literal: "@", Line: 1, Column: 1, EndLine: 1, EndColumn: 2, StartOffset: 0, EndOffset: 1},
		},
		{
			name:  "eof",
			input: "",
			want:  token.Token{Type: token.EOF, Literal: "\000", Line: 1, Column: 1, EndLine: 1, EndColumn: 1, StartOffset: 0, EndOffset: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.NewLexer(tt.input)
			got := lexer.NextToken(l)
			if got.Type != tt.want.Type ||
				got.Literal != tt.want.Literal ||
				got.Line != tt.want.Line ||
				got.Column != tt.want.Column ||
				got.EndLine != tt.want.EndLine ||
				got.EndColumn != tt.want.EndColumn ||
				got.StartOffset != tt.want.StartOffset ||
				got.EndOffset != tt.want.EndOffset {
				t.Errorf("NextToken() = %v, want %v", got, tt.want)
			}
			if got.StartOffset > got.EndOffset {
				t.Fatalf("invalid offsets: start=%d end=%d", got.StartOffset, got.EndOffset)
			}
		})
	}
}
