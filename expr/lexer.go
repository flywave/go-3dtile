package expr

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	TT_INVALID TokenType = iota
	TT_NUMBER
	TT_STRING
	TT_IDENTIFIER
	TT_VARIABLE_OPEN  // ${
	TT_VARIABLE_CLOSE // }
	TT_EOF
	TT_TRUE
	TT_FALSE
	TT_NULL
	TT_UNDEFINED

	// Operators
	TT_PLUS
	TT_MINUS
	TT_STAR
	TT_SLASH
	TT_PERCENT
	TT_EQ     // ===
	TT_NEQ    // !==
	TT_GT     // >
	TT_GTE    // >=
	TT_LT     // <
	TT_LTE    // <=
	TT_AND    // &&
	TT_OR     // ||
	TT_BANG   // !
	TT_TILDE  // ~
	TT_BAR    // |
	TT_REGEX_MATCH // =~
	TT_REGEX_NOT   // !~

	// Delimiters
	TT_DOT
	TT_COMMA
	TT_QUESTION
	TT_COLON
	TT_LPAREN
	TT_RPAREN
	TT_LBRACKET
	TT_RBRACKET

	// Special tokens from lexer
	TT_REGEX // /pattern/flags
)

type Token struct {
	Type  TokenType
	Value string
	Pos   int
}

type Lexer struct {
	input  string
	pos    int
	start  int
	tokens []Token
	err    error
}

func Lex(input string) ([]Token, error) {
	l := &Lexer{input: input}
	l.tokenize()
	return l.tokens, l.err
}

func (l *Lexer) tokenize() {
	for l.pos < len(l.input) {
		l.start = l.pos
		c := l.next()
		if c == 0 {
			break
		}
		if unicode.IsSpace(c) {
			for unicode.IsSpace(l.peek()) {
				l.next()
			}
			l.start = l.pos
			continue
		}
		if c == '/' && l.peek() == '/' {
			for l.peek() != '\n' && l.peek() != 0 {
				l.next()
			}
			l.start = l.pos
			continue
		}
		if c == '/' && l.peek() == '*' {
			l.next()
			for {
				cc := l.next()
				if cc == 0 {
					l.error("unterminated block comment")
					return
				}
				if cc == '*' && l.peek() == '/' {
					l.next()
					break
				}
			}
			l.start = l.pos
			continue
		}
		if c == '/' && l.canStartRegex() {
			l.lexRegex()
			continue
		}
		l.lexToken(c)
		if l.err != nil {
			return
		}
	}
	l.tokens = append(l.tokens, Token{Type: TT_EOF, Pos: l.pos})
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, sz := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += sz
	return r
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *Lexer) emit(tt TokenType) {
	l.tokens = append(l.tokens, Token{Type: tt, Value: l.input[l.start:l.pos], Pos: l.start})
}

func (l *Lexer) emitValue(tt TokenType, val string) {
	l.tokens = append(l.tokens, Token{Type: tt, Value: val, Pos: l.start})
}

func (l *Lexer) error(msg string) {
	if l.err == nil {
		l.err = fmt.Errorf("position %d: %s", l.start, msg)
	}
}

func (l *Lexer) canStartRegex() bool {
	if len(l.tokens) == 0 {
		return true
	}
	last := l.tokens[len(l.tokens)-1]
	switch last.Type {
	case TT_LPAREN, TT_LBRACKET, TT_COMMA, TT_EQ, TT_NEQ, TT_GT, TT_GTE, TT_LT, TT_LTE, TT_AND, TT_OR, TT_BANG, TT_TILDE, TT_PLUS, TT_MINUS, TT_STAR, TT_SLASH, TT_PERCENT, TT_QUESTION, TT_COLON, TT_BAR, TT_VARIABLE_OPEN, TT_REGEX_MATCH, TT_REGEX_NOT:
		return true
	}
	return false
}

func (l *Lexer) lexRegex() {
	l.pos-- // back to the opening /
	l.start = l.pos
	l.next() // consume /
	escape := false
	for {
		c := l.next()
		if c == 0 {
			l.error("unterminated regex")
			return
		}
		if escape {
			escape = false
			continue
		}
		if c == '\\' {
			escape = true
			continue
		}
		if c == '/' {
			break
		}
	}
	for l.peek() == 'g' || l.peek() == 'i' || l.peek() == 'm' {
		l.next()
	}
	l.emit(TT_REGEX)
}

func (l *Lexer) lexToken(c rune) {
	switch {
	case c >= '0' && c <= '9':
		l.lexNumber()

	case c == '\'' || c == '"':
		l.lexString(c)

	case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_':
		l.lexIdent()

	case c == '$' && l.peek() == '{':
		l.next()
		l.emit(TT_VARIABLE_OPEN)

	case c == '}':
		l.emit(TT_VARIABLE_CLOSE)

	default:
		l.lexSymbol(c)
	}
}

func (l *Lexer) lexNumber() {
	start := l.pos - 1
	hasDot := false
	for {
		c := l.peek()
		if c == '.' && !hasDot {
			hasDot = true
			l.next()
			continue
		}
		if c == 'e' || c == 'E' {
			l.next()
			if l.peek() == '+' || l.peek() == '-' {
				l.next()
			}
			for l.peek() >= '0' && l.peek() <= '9' {
				l.next()
			}
			break
		}
		if c >= '0' && c <= '9' {
			l.next()
			continue
		}
		break
	}
	l.emitValue(TT_NUMBER, l.input[start:l.pos])
}

func (l *Lexer) lexString(quote rune) {
	var buf strings.Builder
	for {
		c := l.next()
		if c == 0 {
			l.error("unterminated string")
			return
		}
		if c == '\\' {
			n := l.next()
			switch n {
			case 'n':
				buf.WriteByte('\n')
			case 't':
				buf.WriteByte('\t')
			case 'r':
				buf.WriteByte('\r')
			case '"', '\'', '\\', '/':
				buf.WriteRune(n)
			case 'u':
				hex := l.input[l.pos : l.pos+4]
				if len(hex) >= 4 {
					var val rune
					fmt.Sscanf(hex, "%04x", &val)
					buf.WriteRune(val)
					l.pos += 4
				}
			default:
				buf.WriteRune(n)
			}
			continue
		}
		if c == quote {
			break
		}
		buf.WriteRune(c)
	}
	l.emitValue(TT_STRING, buf.String())
}

func (l *Lexer) lexIdent() {
	start := l.pos - 1
	for {
		c := l.peek()
		if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_' {
			l.next()
			continue
		}
		break
	}
	word := l.input[start:l.pos]
	switch word {
	case "true":
		l.emit(TT_TRUE)
	case "false":
		l.emit(TT_FALSE)
	case "null":
		l.emit(TT_NULL)
	case "undefined":
		l.emit(TT_UNDEFINED)
	default:
		l.emitValue(TT_IDENTIFIER, word)
	}
}

func (l *Lexer) lexSymbol(c rune) {
	switch c {
	case '+':
		l.emit(TT_PLUS)
	case '-':
		l.emit(TT_MINUS)
	case '*':
		l.emit(TT_STAR)
	case '/':
		l.emit(TT_SLASH)
	case '%':
		l.emit(TT_PERCENT)
	case '.':
		l.emit(TT_DOT)
	case ',':
		l.emit(TT_COMMA)
	case '?':
		l.emit(TT_QUESTION)
	case ':':
		l.emit(TT_COLON)
	case '(':
		l.emit(TT_LPAREN)
	case ')':
		l.emit(TT_RPAREN)
	case '[':
		l.emit(TT_LBRACKET)
	case ']':
		l.emit(TT_RBRACKET)
	case '~':
		l.emit(TT_TILDE)
	case '|':
		if l.peek() == '|' {
			l.next()
			l.emit(TT_OR)
		} else {
			l.emit(TT_BAR)
		}
	case '&':
		if l.peek() == '&' {
			l.next()
			l.emit(TT_AND)
		} else {
			l.error("unexpected '&', did you mean &&?")
		}
	case '!':
		if l.peek() == '=' {
			l.next()
			if l.peek() == '=' {
				l.next()
				l.emit(TT_NEQ)
			} else {
				l.error("unexpected '!=', did you mean !==?")
			}
		} else if l.peek() == '~' {
			l.next()
			l.emit(TT_REGEX_NOT)
		} else {
			l.emit(TT_BANG)
		}
	case '=':
		if l.peek() == '=' {
			l.next()
			if l.peek() == '=' {
				l.next()
				l.emit(TT_EQ)
			} else {
				l.error("unexpected '==', did you mean ===?")
			}
		} else if l.peek() == '~' {
			l.next()
			l.emit(TT_REGEX_MATCH)
		} else {
			l.error("unexpected '=', did you mean === or =~?")
		}
	case '>':
		if l.peek() == '=' {
			l.next()
			l.emit(TT_GTE)
		} else {
			l.emit(TT_GT)
		}
	case '<':
		if l.peek() == '=' {
			l.next()
			l.emit(TT_LTE)
		} else {
			l.emit(TT_LT)
		}
	default:
		l.error(fmt.Sprintf("unexpected character: %c", c))
	}
}
