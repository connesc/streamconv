package joiners

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type varintJoinerCommand struct {
	name string
}

func (c *varintJoinerCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *varintJoinerCommand) Parse(args []string, out io.Writer) (joiner streamconv.Joiner, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewVarintJoiner(out), nil
}

func NewVarintJoinerCommand(name string) (command streamconv.JoinerCommand) {
	return &varintJoinerCommand{name}
}

func RegisterVarintJoiner(name string) {
	streamconv.RegisterJoiner(name, NewVarintJoinerCommand(name))
}
