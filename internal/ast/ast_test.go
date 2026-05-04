package ast

import (
	"testing"

	"make_a_lang/internal/lexer"
)

func TestProgramString(t *testing.T) {
	l := lexer.NewLexer("let myVar = anotherVar;")
	p := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: lexer.NextToken(l),
				Name: &Identifier{
					Token: lexer.NextToken(l),
					Value: "myVar",
				},
				Value: &Identifier{
					Token: lexer.NextToken(l),
					Value: "anotherVar",
				},
			},
		},
	}

	want := "let myVar = anotherVar;" 
	if p.String() != want {
		t.Errorf("Program.String() = %q, want %q", p.String(), want)
	}
}
