infix expression:

(114 * 514)

(114 * (514 + (1919 * 810)))

1. Left=114，Operator= *，Right=(514 + (1919 * 810))

2. Left=514，Operator= +，Right=(1919 * 810)


func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // skip ','
		p.nextToken() // move to next identifier

		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}
	
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}