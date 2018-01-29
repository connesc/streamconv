package parser

import (
	"fmt"
	"io"
)

type Program []*Command

type Command struct {
	Words      []string
	SubProgram Program
}

type parser struct {
	lexer *lexer
	ready bool
	token *token
}

func (p *parser) readProgram() (program Program) {
	if !p.ready {
		p.nextToken()
		p.ready = true
	}

	command := &Command{}
	hasSubProgram := false

	for p.token != nil && p.token.kind != tokenRightBrace {
		switch p.token.kind {
		case tokenString:
			if hasSubProgram {
				panic(fmt.Errorf("unexpected token: %v", p.token.value))
			}
			command.Words = append(command.Words, p.token.value)
		case tokenLeftBrace:
			if hasSubProgram {
				panic(fmt.Errorf("unexpected sub-program"))
			}
			p.nextToken()
			command.SubProgram = p.readProgram()
			hasSubProgram = true
			if p.token == nil || p.token.kind != tokenRightBrace {
				panic(fmt.Errorf("closing brace not found"))
			}
		case tokenPipe:
			program = append(program, command)
			command = &Command{}
			hasSubProgram = false
		}

		p.nextToken()
	}

	if len(program) > 0 || len(command.Words) > 0 || len(command.SubProgram) > 0 {
		program = append(program, command)
	}

	return
}

func (p *parser) nextToken() {
	p.token = p.lexer.nextToken()
}

func Parse(reader io.Reader) (program Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			if caughtErr, ok := r.(error); ok {
				err = caughtErr
			} else {
				panic(r)
			}
		}
	}()

	parser := parser{lexer: newLexer(reader)}

	program = parser.readProgram()
	if parser.token != nil {
		err = fmt.Errorf("unexpected token after the end of program")
	}
	return
}
