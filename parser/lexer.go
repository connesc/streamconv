package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"unicode"
)

type tokenKind int

const (
	tokenPipe tokenKind = iota
	tokenLeftBrace
	tokenRightBrace
	tokenString
)

type token struct {
	kind  tokenKind
	value string
}

var escapedChars = map[rune]rune{
	'b': '\b',
	'f': '\f',
	'n': '\n',
	'r': '\r',
	't': '\t',
}

var specialChars = map[rune]*token{
	'|': &token{kind: tokenPipe},
	'{': &token{kind: tokenLeftBrace},
	'}': &token{kind: tokenRightBrace},
}

func isSpecial(char rune) bool {
	_, ok := specialChars[char]
	return ok
}

type lexer struct {
	r     *bufio.Reader
	ready bool
	char  rune
	isEOF bool
}

func newLexer(reader io.Reader) *lexer {
	return &lexer{r: bufio.NewReader(reader)}
}

func (l *lexer) nextToken() (token *token) {
	if !l.ready {
		l.nextChar()
		l.ready = true
	}

	l.skipSpaces()
	if l.isEOF {
		return nil
	}

	if token, found := specialChars[l.char]; found {
		l.nextChar()
		return token
	}

	return l.readString()
}

func (l *lexer) skipSpaces() {
	for !l.isEOF && unicode.IsSpace(l.char) {
		l.nextChar()
	}
}

func (l *lexer) nextChar() {
	char, _, err := l.r.ReadRune()
	if err == nil {
		l.char = char
	} else if err == io.EOF {
		l.isEOF = true
	} else {
		panic(err)
	}
}

func (l *lexer) readString() *token {
	buffer := &bytes.Buffer{}

	escaped := false

	quoted := false
	var quote rune

	complete := false

	for {
		switch {
		case l.isEOF:
			if escaped {
				panic(fmt.Errorf("incomplete escape sequence"))
			}
			if quoted {
				panic(fmt.Errorf("incomplete string (quote not found: %v)", string(quote)))
			}
			complete = true
		case escaped:
			if escapedChar, found := escapedChars[l.char]; found {
				buffer.WriteRune(escapedChar)
			} else {
				buffer.WriteRune(l.char)
			}
			escaped = false
		case l.char == '\\':
			escaped = true
		case quoted:
			if l.char == quote {
				quoted = false
			} else {
				buffer.WriteRune(l.char)
			}
		case isSpecial(l.char) || unicode.IsSpace(l.char):
			complete = true
		case l.char == '\'' || l.char == '"':
			quoted = true
			quote = l.char
		default:
			buffer.WriteRune(l.char)
		}

		if complete {
			return &token{tokenString, buffer.String()}
		}

		l.nextChar()
	}
}
