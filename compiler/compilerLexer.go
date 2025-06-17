/*
- this is a copy of the language/lexer.go, I have just tried to optimize it
- I am going to be doing this for every file under language package
- This will be done until the compiler is fully implemented and then the interpreter
will be removed.
- This is **NOT** completely confirmed yet, but I still plan to make the interpreter into
a JIT compiler. And I also plan to make compilation the default
*/
package compiler

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	STRING
	IDENT
	NEWLINE
	NUM
	COMMENT

	DEF
	ASSIGN
	EQUAL
	MINUS
	ASTERISK
	PLUS
	MODULO
	SLASH
	MORETHAN
	LESSTHAN
	LESSOREQ
	GREATOREQ
	NOTEQ
	NOT
	AND
	OR
	SHELL
	CONCAT

	COMMA
	SEMICOLON
	LPAREN
	RPAREN
	LBRACE
	RBRACE
	LBRACKET
	RBRACKET
	PIPE

	TASK
	EXEC
	WHILE // Imeplent this too
	IF
	ELSE
	IMPORT // TOOD: implement multiple files support
	FOREACH
	COMPILE
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type Lexer struct {
	input    string
	pos      int
	readPos  int
	ch       rune
	line     int
	column   int
	keywords map[string]TokenType
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.keywords = map[string]TokenType{
		"task":    TASK,
		"compile": COMPILE,
		"if":      IF,
		"else":    ELSE,
		"while":   WHILE,
		"foreach": FOREACH,
	}
	return l
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // EOF
	} else {
		l.ch = rune(l.input[l.readPos])
	}

	l.pos = l.readPos
	l.readPos++
	if l.ch == '\n' {
		l.line++
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPos])
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipSpace()
	tok.Line = l.line
	tok.Column = l.column
	switch l.ch {
	case '#':
		tok.Type = COMMENT
		tok.Literal = l.readComment()
		return tok
	case '"', '\'':
		tok.Type = STRING
		tok.Literal = l.readString(l.ch)
		return tok

	// Opertators.
	case ':':
		tok.Type = DEF
		tok.Literal = "="
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			tok.Type = CONCAT
			tok.Literal = "++"
		} else {
			tok.Type = PLUS
			tok.Literal = "+"
		}
	case '-':
		tok.Type = MINUS
		tok.Literal = "-"
	case '/':
		tok.Type = SLASH
		tok.Literal = "/"
	case '*':
		tok.Type = ASTERISK
		tok.Literal = "*"
	case '%':
		tok.Type = MODULO
		tok.Literal = "%"
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = EQUAL
			tok.Literal = "=="
		} else {
			tok.Type = ASSIGN
			tok.Literal = "="

		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok.Type = NOTEQUAL
			tok.Literal = "!="
		} else {
			tok.Type = NOT
			tok.Literal = "!"
		}
	case '<':
		if l.peekChar() == '=' {
			tok.Type = LORETO
			tok.Literal = "<="
		} else {
			tok.Type = LESSTHAN
			tok.Literal = "<"
		}
	case '>':
		if l.peekChar() == '=' {
			tok.Type = GORETO
			tok.Literal = ">="
		} else {
			tok.Type = GREATERTHAN
			tok.Literal = ">"
		}
	case '|':
		if l.peekChar() == '|' {
			tok.Type = OR
			tok.Literal = "||"
		} else {
			tok.Type = PIPE
			tok.Literal = "|"
		}
	case '&':
		if l.peekChar() == '&' {
			tok.Type = AND
			tok.Literal = "&&"
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch) // idk why string() this time lol
		}
	case '$':
		tok.Type = SHELL
		tok.Literal = "$"

	// Delimiters.
	case '(':
		tok.Type = LPAREN
		tok.Literal = "("
	case ')':
		tok.Type = RPAREN
		tok.Literal = ")"
	case '{':
		tok.Type = LBRACE
		tok.Literal = "{"

	case '}':
		tok.Type = RBRACE
		tok.Literal = "}"

	case '[':
		tok.Type = LBRACKET
		tok.Literal = "["

	case ']':
		tok.Type = RBRACKET
		tok.Literal = "]"
	case ',':
		tok.Type = COMMA
		tok.Literal = ","
	case ';':
		tok.Type = SEMICOLON
		tok.Literal = ";"

	// Newline and EOF.
	case '\n':
		tok.Type = NEWLINE
		tok.Literal = "\n"

	case 0:
		tok.Type = EOF
		tok.Literal = ""

	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			if tokenType, ok := l.keywords[tok.Literal]; ok {
				tok.Type = tokenType
			} else {
				tok.Type = IDENT
			}
			return tok
		} else if isDigit(l.ch) {
			tok.Type = NUM
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok.Type = ILLEGAL
			tok.Literal = string(l.ch)
		}

	}

	l.readChar()
	return tok
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func (l *Lexer) skipSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readComment() string {
	pos := l.pos
	for l.ch != '\n' && l.ch != '0' {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}

	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[pos:l.pos]
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}
