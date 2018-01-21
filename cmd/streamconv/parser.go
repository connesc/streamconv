package main

import (
	"bytes"
	"unicode"
)

func parse(program string) (commands [][]string, err error) {
	commands = make([][]string, 0)
	tokens := make([]string, 0)
	token := &bytes.Buffer{}

	escaped := false

	quoted := false
	var quote rune

	endToken := func() {
		if token.Len() > 0 {
			tokens = append(tokens, string(token.Bytes()))
			token.Reset()
		}
	}

	endCommand := func() {
		endToken()
		commands = append(commands, tokens)
		tokens = make([]string, 0)
	}

	for _, char := range program {
		switch {
		case escaped:
			switch char {
			case 'b':
				char = '\b'
			case 'f':
				char = '\f'
			case 'n':
				char = '\n'
			case 'r':
				char = '\r'
			case 't':
				char = '\t'
			}
			escaped = false
			_, err = token.WriteRune(char)
			if err != nil {
				return
			}
		case char == '\\':
			escaped = true
		case quoted:
			if char == quote {
				quoted = false
			} else {
				_, err = token.WriteRune(char)
				if err != nil {
					return
				}
			}
		case char == '|':
			endCommand()
		case unicode.IsSpace(char):
			endToken()
		case char == '\'' || char == '"':
			quoted = true
			quote = char
		default:
			_, err = token.WriteRune(char)
			if err != nil {
				return
			}
		}
	}

	endCommand()
	return
}
