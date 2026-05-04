package ast

import (
	"bytes"
	"fmt"

	token "make_a_lang/internal/token"
)

type Node interface {
	TokenLiteral() string
	Pos() token.Token // 返回节点起始 token（含位置信息）
	End() token.Token // 返回节点结束 token（含位置信息，End* 为尾后）
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) Pos() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[0].Pos()
	}
	return token.Token{}
}

func (p *Program) End() token.Token {
	if len(p.Statements) > 0 {
		return p.Statements[len(p.Statements)-1].End()
	}
	return token.Token{}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// ----- Expressions -----

// 变量名
type Identifier struct {
	Token token.Token
	Value string
}

func (e *Identifier) expressionNode()      {}
func (e *Identifier) TokenLiteral() string { return e.Token.Literal }
func (e *Identifier) Pos() token.Token     { return e.Token }
func (e *Identifier) End() token.Token     { return e.Token }
func (e *Identifier) String() string       { return e.Value }

type IfExpression struct {
	Token       token.Token // 'if' token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement // else部分，可能为nil
}

func (e *IfExpression) expressionNode()      {}
func (e *IfExpression) TokenLiteral() string { return e.Token.Literal }
func (e *IfExpression) Pos() token.Token     { return e.Token }
func (e *IfExpression) End() token.Token {
	if e.Alternative != nil {
		return e.Alternative.End()
	}
	return e.Consequence.End()
}
func (e *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if (")
	out.WriteString(e.Condition.String())
	out.WriteString(") ")
	out.WriteString(e.Consequence.String())
	if e.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(e.Alternative.String())
	}
	return out.String()
}

// 前缀表达式
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (e *PrefixExpression) expressionNode()      {}
func (e *PrefixExpression) TokenLiteral() string { return e.Token.Literal }
func (e *PrefixExpression) Pos() token.Token     { return e.Token }
func (e *PrefixExpression) End() token.Token     { return e.Right.End() }
func (e *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right.String())
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (e *InfixExpression) expressionNode()      {}
func (e *InfixExpression) TokenLiteral() string { return e.Token.Literal }
func (e *InfixExpression) Pos() token.Token     { return e.Left.Pos() }
func (e *InfixExpression) End() token.Token     { return e.Right.End() }
func (e *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left.String(), e.Operator, e.Right.String())
}

type CallExpression struct {
	Token    token.Token
	Function Expression   // 可能是Identifier，也可能是FunctionLiteral
	Arguments []Expression
}

func (e *CallExpression) expressionNode()      {}
func (e *CallExpression) TokenLiteral() string { return e.Token.Literal }
func (e *CallExpression) Pos() token.Token     { return e.Function.Pos() }
func (e *CallExpression) End() token.Token {
	if len(e.Arguments) > 0 {
		return e.Arguments[len(e.Arguments)-1].End()
	}
	return e.Function.End()
}
func (e *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(e.Function.String())
	out.WriteString("(")
	for i, arg := range e.Arguments {
		out.WriteString(arg.String())
		if i < len(e.Arguments)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(")")
	return out.String()
}

type IndexExpression struct {
	Token token.Token // '[' token
	Left  Expression  // array
	Index Expression  // index
}

func (e *IndexExpression) expressionNode()      {}
func (e *IndexExpression) TokenLiteral() string { return e.Token.Literal }
func (e *IndexExpression) Pos() token.Token     { return e.Left.Pos() }
func (e *IndexExpression) End() token.Token     { return e.Index.End() }
func (e *IndexExpression) String() string {
	return fmt.Sprintf("(%s[%s])", e.Left.String(), e.Index.String())
}

// ----- Statements -----

type LetStatement struct {
	Token token.Token // LET token
	Name  *Identifier // let变量的变量名
	Value Expression  // let变量的值
}

func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal } // 返回 "let"
func (s *LetStatement) Pos() token.Token     { return s.Token }
func (s *LetStatement) End() token.Token {
	if s.Value != nil {
		return s.Value.End()
	}
	return s.Token
}

func (s *LetStatement) String() string {
	var out bytes.Buffer
	// first "let"
	out.WriteString(s.TokenLiteral())
	out.WriteString(" ")
	if s.Name != nil {
		out.WriteString(s.Name.String())
	}
	out.WriteString(" = ")
	if s.Value != nil {
		out.WriteString(s.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type AssignStatement struct {
	Token token.Token // ASSIGN
	Name  *Identifier
	Value Expression
}

func (s *AssignStatement) statementNode()       {}
func (s *AssignStatement) TokenLiteral() string { return s.Token.Literal } // 返回 "="
func (s *AssignStatement) Pos() token.Token {
	if s.Name != nil {
		return s.Name.Token
	}
	return s.Token
}

func (s *AssignStatement) End() token.Token {
	if s.Value != nil {
		return s.Value.End()
	}
	return s.Token
}

func (s *AssignStatement) String() string {
	var out bytes.Buffer
	if s.Name != nil {
		out.WriteString(s.Name.String())
	}
	out.WriteString(" = ")
	if s.Value != nil {
		out.WriteString(s.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token // RETURN
	ReturnValue Expression  // return后面的表达式
}

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal } // 返回 "return"
func (s *ReturnStatement) Pos() token.Token     { return s.Token }
func (s *ReturnStatement) End() token.Token {
	if s.ReturnValue != nil {
		return s.ReturnValue.End()
	}
	return s.Token
}

func (s *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(s.TokenLiteral())
	out.WriteString(" ")
	if s.ReturnValue != nil {
		out.WriteString(s.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token // 表达式首 token
	Expression Expression
}

func (s *ExpressionStatement) statementNode()       {}
func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ExpressionStatement) Pos() token.Token     { return s.Token }
func (s *ExpressionStatement) End() token.Token {
	if s.Expression != nil {
		return s.Expression.End()
	}
	return s.Token
}

func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}
	return ""
}

type BlockStatement struct {
	Token      token.Token // LBRACE
	Statements []Statement
}

func (s *BlockStatement) statementNode()       {}
func (s *BlockStatement) TokenLiteral() string { return s.Token.Literal }
func (s *BlockStatement) Pos() token.Token     { return s.Token }
func (s *BlockStatement) End() token.Token {
	if len(s.Statements) > 0 {
		return s.Statements[len(s.Statements)-1].End()
	}
	return s.Token
}

func (s *BlockStatement) String() string {
	var out bytes.Buffer
	for _, stmt := range s.Statements {
		out.WriteString(stmt.String())
	}
	return out.String()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (e *IntegerLiteral) expressionNode()      {}
func (e *IntegerLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *IntegerLiteral) Pos() token.Token     { return e.Token }
func (e *IntegerLiteral) End() token.Token     { return e.Token }
func (e *IntegerLiteral) String() string       { return e.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (e *BooleanLiteral) expressionNode()      {}
func (e *BooleanLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *BooleanLiteral) Pos() token.Token     { return e.Token }
func (e *BooleanLiteral) End() token.Token     { return e.Token }
func (e *BooleanLiteral) String() string       { return e.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (e *StringLiteral) expressionNode()      {}
func (e *StringLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *StringLiteral) Pos() token.Token     { return e.Token }
func (e *StringLiteral) End() token.Token     { return e.Token }
func (e *StringLiteral) String() string       { return e.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token // 'fun' token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (e *FunctionLiteral) expressionNode()      {}
func (e *FunctionLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *FunctionLiteral) Pos() token.Token     { return e.Token }
func (e *FunctionLiteral) End() token.Token {
	if e.Body != nil {
		return e.Body.End()
	}
	return e.Token
}
func (e *FunctionLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(e.TokenLiteral())
	out.WriteString("(")
	for i, param := range e.Parameters {
		out.WriteString(param.String())
		if i < len(e.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	if e.Body != nil {
		out.WriteString(e.Body.String())
	}
	return out.String()
}