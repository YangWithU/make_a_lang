package lexer

type Lexer struct {
	input string
	position int
	readPosition int // always == position + 1
	ch byte // current character
	line int
	column int
}

// core function to read the next character and advance the positions
func (l *Lexer) readChar() {
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}

	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII code for NUL, signifies end of file
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
	l.column++
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

func NextToken(l *Lexer) Token {
	skipWhitespace(l)
	skipComment(l)
	// TODO
}

func NewLexer(input string) *Lexer {
	l := Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return &l // escapes to heap
}