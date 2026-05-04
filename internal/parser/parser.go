package parser

import (
	"fmt"
	"strconv"

	"make_a_lang/internal/ast"
	"make_a_lang/internal/lexer"
	"make_a_lang/internal/token"
)

type Precedence int

const (
	_ Precedence = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

var precedenceMap = map[token.TokenType]Precedence{
	token.EQ:     EQUALS,
	token.NOT_EQ: EQUALS,

	token.LT: LESSGREATER,
	token.GT: LESSGREATER,

	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,

	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,

	token.ASSIGN: LOWEST,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type ParseError struct {
	Message string
	Token   token.Token
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token
	peekToken token.Token

	errors []ParseError

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

func (p *Parser) newError(msg string, tok token.Token) {
	p.errors = append(p.errors, ParseError{Message: msg, Token: tok})
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = lexer.NextToken(p.l)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         []ParseError{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)         // add, foobar, x, y
	p.registerPrefix(token.INT, p.parseIntegerLiteral)       // 1343456
	p.registerPrefix(token.STRING, p.parseStringLiteral)     // "foobar"
	p.registerPrefix(token.IF, p.parseIfExpression)          // if (condition) { consequence } else { alternative }
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral) // fun(x) { ... }
	p.registerPrefix(token.BANG, p.parsePrefixExpression)    // !true
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)   // -1
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression) // (1 + 2)
	p.registerPrefix(token.TRUE, p.parseBoolean)            // true
	p.registerPrefix(token.FALSE, p.parseBoolean)           // false

	// infix operators
	p.registerInfix(token.PLUS, p.parseInfixExpression)     // 1 + 2
	p.registerInfix(token.MINUS, p.parseInfixExpression)    // 1 - 2
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) // 1 * 2
	p.registerInfix(token.SLASH, p.parseInfixExpression)    // 1 / 2
	p.registerInfix(token.EQ, p.parseInfixExpression)       // 1 == 2
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)   // 1 != 2
	p.registerInfix(token.LT, p.parseInfixExpression)       // 1 < 2
	p.registerInfix(token.GT, p.parseInfixExpression)       // 1 > 2

	p.registerInfix(token.LPAREN, p.parseCallExpression)    // add(1, 2)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression) // myArray[1]

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expect nextToken type: [%v], actual: [%v]", p.peekToken, t)
	p.newError(msg, p.peekToken)
}

// if peekToken.type == t then headon to nextToken
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENT, token.INT, token.BANG, token.MINUS:
		return p.parseExpressionStatement()
	default:
		if p.curTokenIs(token.IDENT) && p.peekTokenIs(token.ASSIGN) {
			return p.parseAssignStatement()
		} else {
			return p.parseExpressionStatement()
		}
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	root := &ast.Program{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			root.Statements = append(root.Statements, stmt)
		}
		p.nextToken()
	}
	return root
}

func (p *Parser) curPrecedence() Precedence {
	if p, ok := precedenceMap[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() Precedence {
	if p, ok := precedenceMap[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// ----- identifiers, literals -----

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.newError(fmt.Sprintf("could not parse %q as integer", p.curToken.Literal), p.curToken)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	identifiers := []*ast.Identifier{}
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	p.nextToken()

	for p.curTokenIs(token.COMMA) {
		p.nextToken() // ,
		p.nextToken() // ident

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return identifiers
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	params := p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	body := p.parseBlockStatement()

	return &ast.FunctionLiteral{
		Token:      p.curToken,
		Parameters: params,
		Body:       body,
	}
}

// ----- expressions -----

// 给当前表达式优先级，如果peekToken的优先级更高，则继续解析peekToken为当前表达式的右侧表达式
func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefixFn := p.prefixParseFns[p.curToken.Type]
	if prefixFn == nil {
		p.newError(fmt.Sprintf("no prefix parse method for [%v]", p.curToken.Type), p.curToken)
		return nil
	}

	leftExpr := prefixFn()

	// 解析例如 1+2*3
	// 发现后面还有更高优先级的运算符，就继续解析右侧表达式
	if p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infixFn := p.infixParseFns[p.peekToken.Type]
		if infixFn == nil {
			return leftExpr
		}

		p.nextToken() // 先把peekToken移到curToken上，再调用infixFn解析右侧表达式
		leftExpr = infixFn(leftExpr)
	}
	return leftExpr
}

// -5, !true
func (p *Parser) parsePrefixExpression() ast.Expression {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Right:    nil,
	}

	p.nextToken() // now transit, have numbers and exprs

	expr.Right = p.parseExpression(PREFIX) // recursive call parseExpression
	// with higher precedence to parse the right side of the prefix expression
	return expr
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expr := ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
		Right:    nil,
	}

	prec := p.curPrecedence()
	p.nextToken() // now transit, have numbers and exprs

	expr.Right = p.parseExpression(prec)
	return &expr
}

// add(1, 2), myArray[3], fun(x)
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{
		Token:     p.curToken,
		Function:  function,
		Arguments: []ast.Expression{},
	}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken() // skip ')'
		return expr
	}

	p.nextToken() // skip '('

	expr.Arguments = append(expr.Arguments, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip ','
		p.nextToken() // move to next argument
		expr.Arguments = append(expr.Arguments, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	expr := &ast.IndexExpression{
		Token: p.curToken,
		Left:  left,
		Index: nil,
	}

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken() // skip ']'
		return expr
	}

	p.nextToken() // skip '['

	expr.Index = p.parseExpression(LOWEST)

	if expr.Index == nil {
		return nil
	}

	if expr.Index.Pos().Type != token.INT {
		p.newError(fmt.Sprintf("index expression index must be integer, got %v", expr.Index.Pos().Type), expr.Index.Pos())
		return nil
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return expr
}

// ----- statements -----

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := ast.ReturnStatement{Token: p.curToken}

	p.nextToken() // skip 'return' keyword

	for !p.curTokenIs(token.SEMICOLON) && !p.curTokenIs(token.EOF) {
		p.nextToken()
	}

	return &stmt
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
		Name:  nil,
		Value: nil,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken() // skip '='

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseAssignStatement() ast.Statement {
	stmt := &ast.AssignStatement{
		Token: p.peekToken, // ASSIGN token
		Name:  nil,
		Value: nil,
	}

	if !p.curTokenIs(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken() // skip '='

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// ((expr))
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // skip '('

	exp := p.parseExpression(LOWEST) // start a new expression with lowest precedence

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.curToken,
		Statements: []ast.Statement{},
	}

	p.nextToken() // skip '{'

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{
		Token: p.curToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken() // skip '('
	expr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // skip 'else'

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}
