package language

type TokenType int

const (
	// Token types.
	EOF TokenType = iota
	ILLEGAL

	// Literals.
	STRING
	IDENT
	NEWLINE
	NUMBER
	COMMENT

	// Operators.
	DEFINE      // :
	ASSIGN      // =
	EQUAL       // ==
	MINUS       // -
	ASTERISK    // *
	PLUS        // +
	MODULO      // %
	SLASH       // `/`
	GREATERTHAN // >
	LESSTHAN    // <
	LORETO      // <=
	GORETO      // >=
	NOTEQUAL    // !=
	NOT         // !
	AND         // &&
	OR          // ||
	SHELL       // $
	CONCAT      // ++ (new concatenation operator [fire])

	// Delimiters.
	COMMA     // `,`
	SEMICOLON // `;`
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	PIPE      // |

	// Keywords.
	TASK       // task
	RUN        // run (run command)
	IF         // if
	ELSE       // else
	IMPORT     // include
	DEPENDENCY // (require programs) require
	SWAP       // (swap two variables) swap
	WHILE      // while
	FOREACH    // (foreach thing in an array or some shit idk) foreach
	COMPILE    // (compile things with command)  compile
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type Lexer struct {
	input        string
	position     int  // current pos in input
	readPosition int  // current reading pos
	ch           rune // current char under examination
	line         int  // current line
	column       int  // current column num
	keywords     map[string]TokenType
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.keywords = map[string]TokenType{
		"task":     TASK,
		"compile":  COMPILE,
		"require":  DEPENDENCY,
		"if":       IF,
		"else":     ELSE,
		"swap":     SWAP,
		"requires": DEPENDENCY,
		"while":    WHILE,
		"foreach":  FOREACH,
	}

	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // represents EOF
	} else {
		l.ch = rune(l.input[l.readPosition])
	}

	l.position = l.readPosition
	l.readPosition++

	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPosition])
}

func (l *Lexer) NextToken() Token {
	var tok Token
	l.skipWhitespace()
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	// Literals.
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
		tok.Type = DEFINE
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
			tok.Type = NUMBER
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

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	// To handle floating point numbers.
	if l.ch == '.' && isDigit(l.peekChar()) {
		l.readChar()
		for isDigit(l.ch) {
			l.readChar()
		}
	}
	return l.input[position:l.position]
}

func isDigit(thing rune) bool {
	return '0' <= thing && thing <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '-' || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(a rune) bool {
	return 'a' <= a && a <= 'z' || 'A' <= a && a <= 'Z'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readComment() string {
	position := l.position
	for l.ch != '\n' && l.ch != '0' {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readString(quote rune) string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == quote {
			break
		}

		if l.ch == '\\' && l.peekChar() == quote {
			l.readChar()
		}
	}

	if l.ch == 0 {
		return l.input[position-1 : l.position]
	}

	result := l.input[position:l.position]
	l.readChar()
	return result
}
