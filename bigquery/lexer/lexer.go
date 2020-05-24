package lexer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/tadyjp/bigquery-vis/bigquery/token"
)

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input  *bufio.Reader
	buffer bytes.Buffer
	state  stateFn
	pos    int
	start  int
	tokens chan token.Token
}

const eof = -1

func (l *Lexer) NextToken() token.Token {
	return <-l.tokens
}

func New(input io.Reader) *Lexer {
	l := &Lexer{
		input:  bufio.NewReader(input),
		tokens: make(chan token.Token),
	}
	go l.run()
	return l
}

func (l *Lexer) run() {
	for l.state = lexText; l.state != nil; {
		l.state = l.state(l)
	}
}

func (l *Lexer) next() rune {
	r := l.skip()
	if r == eof {
		return eof
	}
	l.buffer.WriteRune(r)
	return r
}

func (l *Lexer) skip() rune {
	r, w, err := l.input.ReadRune()
	if err == io.EOF {
		return eof
	}
	l.pos += w
	return r
}

func (l *Lexer) peek() rune {
	lead, err := l.input.Peek(1)
	if err == io.EOF {
		return eof
	} else if err != nil {
		l.errorf("%s", err.Error())
		return 0
	}

	p, err := l.input.Peek(runeLen(lead[0]))
	if err == io.EOF {
		return eof
	} else if err != nil {
		l.errorf("%s", err.Error())
		return 0
	}
	r, _ := utf8.DecodeRune(p)
	return r
}

func runeLen(lead byte) int {
	if lead < 0xC0 {
		return 1
	} else if lead < 0xE0 {
		return 2
	} else if lead < 0xF0 {
		return 3
	} else {
		return 4
	}
}

func (l *Lexer) emit(t token.TokenType) {
	l.tokens <- token.Token{Type: t, Start: l.start, Literal: l.buffer.String()}
	l.start = l.pos
	l.buffer.Truncate(0)
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.peek()) >= 0 {
		l.next()
		return true
	}
	return false
}

// TODO: accept で置き換え可能？
func (l *Lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.peek()) >= 0 {
		l.next()
	}
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token.Token{Type: token.ILLEGAL, Start: l.start, Literal: fmt.Sprintf(format, args...)}
	return nil
}

func (l *Lexer) hasPrefix(prefix string) bool {
	p, err := l.input.Peek(len(prefix))
	if err == io.EOF {
		return false
	} else if err != nil {
		l.errorf("%s", err.Error())
		return false
	}
	return string(p) == prefix
}

// Accept next count runes. Normally called after hasPrefix().
func (l *Lexer) nextRuneCount(count int) {
	for i := 0; i < count; i++ {
		// l.next()
		l.skip()
	}
}

const (
	doubleQuote  = `"`
	lineComment  = "//"
	leftComment  = "/*"
	rightComment = "*/"
)

func lexText(l *Lexer) stateFn {
Loop:
	for {
		if l.hasPrefix(`"""`) {
			l.nextRuneCount(3)
			return lexString(`"""`)
		}

		r := l.peek()
		switch r {
		case ';':
			l.next()
			l.emit(token.SEMICOLON)
			continue
		case '.':
			l.next()
			l.emit(token.PERIOD)
			continue
		case '*':
			l.next()
			l.emit(token.ASTERISK)
			continue
		case '"':
			l.skip()
			return lexString(`"`)
		case '\'':
			l.skip()
			return lexString(`'`)
		case '`':
			l.skip()
			return lexQuotedIdentifier
			// case ':':
			// 	l.next()
			// 	l.emit(itemColon)
			// case ',':
			// 	l.next()
			// 	l.emit(itemComma)
			// case '.':
			// 	l.next()
			// 	l.emit(itemPeriod)
			// case '{':
			// 	l.next()
			// 	l.emit(itemLeftBrace)
			// case '}':
			// 	l.next()
			// 	l.emit(itemRightBrace)
			// case '[':
			// 	l.next()
			// 	l.emit(itemLeftBracket)
			// case ']':
			// 	l.next()
			// 	l.emit(itemRightBracket)
			// case '/':
			// 	if l.hasPrefix(lineComment) {
			// 		return lexLineComment
			// 	} else if l.hasPrefix(leftComment) {
			// 		return lexBlockComment
			// 	} else {
			// 		return l.errorf("invalid character after slash")
			// 	}
		}

		if unicode.IsSpace(r) {
			l.skip()
			continue
		} else if r == eof {
			l.next()
			break Loop
		} else if strings.IndexRune("0123456789-", r) >= 0 {
			return lexNumber
		} else {
			return lexIdentifier
		}
	}
	l.emit(token.EOF)
	return nil
}

const (
	hexdigit  = "0123456789ABCDEFabcdef"
	digit     = "0123456789"
	digit1To9 = "123456789"
)

func lexString(quote string) func(*Lexer) stateFn {
	return func(l *Lexer) stateFn {
		for {
			if l.hasPrefix(`\` + quote) { // escaped quote
				l.skip()
				l.next()
			}

			r := l.next()
			if r == '\\' {
				if l.accept(`"\/bfnrt`) {
					// break
					// do nothing
				} else if r := l.next(); r == 'u' {
					for i := 0; i < 4; i++ {
						if !l.accept(hexdigit) {
							return l.errorf("expected 4 hexadecimal digits")
						}
					}
				} else {
					return l.errorf("unsupported escape character")
				}
			} else if l.hasPrefix(quote) { // finish on end quote
				l.emit(token.STRING)
				l.nextRuneCount(utf8.RuneCountInString(quote))
				return lexText
			} else if unicode.IsControl(r) {
				return l.errorf("cannot contain control characters in strings")
			} else if r == eof {
				return l.errorf("unclosed string")
			}
		}
	}
}

// TODO: Quoted identifiers cannot be empty.
func lexQuotedIdentifier(l *Lexer) stateFn {
	for {
		r := l.next()

		if l.peek() == '`' {
			l.emit(token.IDENT)
			l.skip()
			return lexText
		} else if r == eof {
			return l.errorf("unclosed quoted string")
		}
	}
}

func lexNumber(l *Lexer) stateFn {
	l.accept("-")
	if l.accept(digit1To9) {
		l.acceptRun(digit)
	} else if !l.accept("0") {
		return l.errorf("bad digit for number")
	}
	if l.accept(".") {
		l.acceptRun(digit)
	}
	if l.accept("eE") {
		l.accept("+-")
		if !l.accept(digit) {
			return l.errorf("digit expected for number exponent")
		}
		l.acceptRun(digit)
	}
	l.emit(token.NUMBER)
	return lexText
}

func lexIdentifier(l *Lexer) stateFn {
	r := l.peek()
	if unicode.IsLetter(r) || r == '$' || r == '_' {
		// do nothing
	} else if r == '\\' {
		if !l.accept("u") {
			return l.errorf("'u' for unicode escape sequence expected")
		}
		for i := 0; i < 4; i++ {
			if !l.accept(hexdigit) {
				return l.errorf("expected 4 hexadecimal digits for unicode escape sequence")
			}
		}
	} else {
		return l.errorf("identifier expected")
	}

	for r = l.peek(); isIdentifierPart(r); {
		l.next()
		r = l.peek()
	}

	literal := l.buffer.String()
	l.emit(token.LookupIdent(strings.ToLower(literal)))
	return lexText
}

func isIdentifierPart(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
}

// func lexLineComment(l *Lexer) stateFn {
// 	for {
// 		r := l.next()
// 		if r == '\n' || r == eof {
// 			if l.pos > l.start {
// 				l.emit(itemLineComment)
// 			}
// 			if r == eof {
// 				l.emit(itemEOF)
// 				return nil
// 			}
// 			return lexText
// 		}
// 	}
// }
//
// func lexBlockComment(l *Lexer) stateFn {
// 	for {
// 		if l.hasPrefix(rightComment) {
// 			l.nextRuneCount(utf8.RuneCountInString(rightComment))
// 			if l.pos > l.start {
// 				l.emit(itemBlockComment)
// 			}
// 			return lexText
// 		}
// 		if l.next() == eof {
// 			break
// 		}
// 	}
// 	l.emit(itemEOF)
// 	return nil
// }
