package joiners

import (
	"fmt"
	"io"
	"streamconv"
)

type simpleJoinerCommand struct {
	name string
}

func (c *simpleJoinerCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *simpleJoinerCommand) Parse(args []string, out io.Writer) (joiner streamconv.Joiner, err error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments (expected up to 1, got %v)", len(args))
	}

	delimiter := "\n"
	if len(args) == 1 {
		delimiter = args[0]
	}

	return NewSimpleJoiner(out, delimiter), nil
}

func NewSimpleJoinerCommand(name string) (command streamconv.JoinerCommand) {
	return &simpleJoinerCommand{name}
}

func RegisterSimpleJoiner(name string) {
	streamconv.RegisterJoiner(name, NewSimpleJoinerCommand(name))
}
