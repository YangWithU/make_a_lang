package lexer

import (
	token "make_a_lang/internal/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int  // always == position + 1
	ch           byte // current character
	line         int  // current line number
	column       int  // current column number

	nextLine   int // the line number for the next character to read
	nextColumn int // the column number for the next character to read
}

type position struct {
	line   int
	column int
	offset int
}

func (l *Lexer) curPos() position {
	return position{line: l.line, column: l.column, offset: l.position}
}

func (l *Lexer) tokenFrom(start position, tokenType token.TokenType, literal string) token.Token {
	end := l.curPos()
	return token.Token{
		Type:        tokenType,
		Literal:     literal,
		Line:        start.line,
		Column:      start.column,
		EndLine:     end.line,
		EndColumn:   end.column,
		StartOffset: start.offset,
		EndOffset:   end.offset,
	}
}

// consume advance number of characters and return a token with the given type and literal, 
// using the start position for token metadata
func (l *Lexer) advanceAndToken(start position, tokenType token.TokenType, literal string, advance int) token.Token {
	for i := 0; i < advance; i++ {
		l.readChar()
	}
	return l.tokenFrom(start, tokenType, literal)
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// core function to read the next character and advance the positions
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
		l.position = l.readPosition
		l.line = l.nextLine
		l.column = l.nextColumn
		return
	}

	ch := l.input[l.readPosition]
	l.ch = ch
	l.position = l.readPosition
	l.line = l.nextLine
	l.column = l.nextColumn

	l.readPosition++

	if l.ch == '\n' {
		l.nextLine++
		l.nextColumn = 1
	} else {
		l.nextColumn++
	}
}

func peekChar(l *Lexer) byte {
	if l.readPosition >= len(l.input) {
		return 0 // EOF
	}
	return l.input[l.readPosition]
}

func skipWhitespace(l *Lexer) {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func skipComment(l *Lexer) {
	if l.ch == '/' {
		if peekChar(l) == '/' {
			// 单行注释，跳过到行尾
			for l.ch != '\n' && l.ch != 0 {
				l.readChar()
			}
		} else if peekChar(l) == '*' {
			// 多行注释，跳过到结束符 */
			l.readChar() // consume '*'
			l.readChar() // move to next char after '/*'
			for {
				if l.ch == 0 {
					break // EOF reached without closing comment
				}
				if l.ch == '*' && peekChar(l) == '/' {
					l.readChar() // consume '*'
					l.readChar() // consume '/'
					break
				}
				l.readChar()
			}
		}
	}
}

func (l *Lexer) readIdentifier() string {
	old_pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[old_pos:l.position]
}

func (l *Lexer) readNumber() string {
	old_pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[old_pos:l.position]
}

func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	old_pos := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	str := l.input[old_pos:l.position]
	l.readChar() // skip closing quote
	return str
}

func NextToken(l *Lexer) token.Token {
	for {
		skipWhitespace(l)
		if l.ch == '/' && (peekChar(l) == '/' || peekChar(l) == '*') {
			skipComment(l)
			continue
		}
		break
	}

	start := l.curPos() // capture the start position for the token
	switch l.ch {
	case '+':
		return l.advanceAndToken(start, token.PLUS, string(l.ch), 1)
	case '-':
		return l.advanceAndToken(start, token.MINUS, string(l.ch), 1)
	case '*':
		return l.advanceAndToken(start, token.ASTERISK, string(l.ch), 1)
	case '/':
		return l.advanceAndToken(start, token.SLASH, string(l.ch), 1)
	case '(':
		return l.advanceAndToken(start, token.LPAREN, string(l.ch), 1)
	case ')':
		return l.advanceAndToken(start, token.RPAREN, string(l.ch), 1)
	case '{':
		return l.advanceAndToken(start, token.LBRACE, string(l.ch), 1)
	case '}':
		return l.advanceAndToken(start, token.RBRACE, string(l.ch), 1)
	case ',':
		return l.advanceAndToken(start, token.COMMA, string(l.ch), 1)
	case ';':
		return l.advanceAndToken(start, token.SEMICOLON, string(l.ch), 1)
	case '<':
		return l.advanceAndToken(start, token.LT, string(l.ch), 1)
	case '>':
		return l.advanceAndToken(start, token.GT, string(l.ch), 1)
	case 0:
		return l.tokenFrom(start, token.EOF, "\000")

	case '=':
		if peekChar(l) == '=' {
			return l.advanceAndToken(start, token.EQ, "==", 2)
		}
		return l.advanceAndToken(start, token.ASSIGN, string(l.ch), 1)

	case '!':
		if peekChar(l) == '=' {
			return l.advanceAndToken(start, token.NOT_EQ, "!=", 2)
		}
		return l.advanceAndToken(start, token.BANG, string(l.ch), 1)

	default:
		// handle string literals, identifiers, and numbers
		if isLetter(l.ch) {
			ident := l.readIdentifier()
			tokenType := token.LookupIdent(ident)
			return l.tokenFrom(start, tokenType, ident)
		} else if isDigit(l.ch) {
			num := l.readNumber()
			return l.tokenFrom(start, token.INT, num)
		} else if l.ch == '"' {
			str := l.readString()
			return l.tokenFrom(start, token.STRING, str)
		}
	}
	return l.advanceAndToken(start, token.ILLEGAL, string(l.ch), 1)
}

func NewLexer(input string) *Lexer {
	l := Lexer{input: input, nextLine: 1, nextColumn: 1}
	l.readChar()
	return &l // escapes to heap
}
