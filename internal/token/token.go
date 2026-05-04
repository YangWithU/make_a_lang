package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string

	Line   int // 起始行（1-based，包含）
	Column int // 起始列（1-based，包含）

	EndLine     int // 终止行（1-based，尾后，不包含）
	EndColumn   int // 终止列（1-based，尾后，不包含）
	StartOffset int // 起始偏移（0-based，包含）
	EndOffset   int // 终止偏移（0-based，尾后，不包含）
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"  // add, foobar, x, y, ...
	INT    = "INT"    // 1343456
	STRING = "STRING" // "foo bar"

	// 运算符
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// 关键字
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	NULL     = "NULL"
)

var keywords = map[string]TokenType{
	"fun":    FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"null":   NULL,
}

// LookupIdent 用于检查给定的标识符是否是关键字
// 如果是关键字，则返回对应的 TokenType；否则返回 IDENT（普通用户标识符）
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
