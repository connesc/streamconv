package combiners

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type varintCombinerCommand struct {
	name string
}

func (c *varintCombinerCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *varintCombinerCommand) Parse(args []string, out io.Writer) (combiner streamconv.ItemWriter, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewVarintCombiner(out), nil
}

func NewVarintCombinerCommand(name string) (command streamconv.CombinerCommand) {
	return &varintCombinerCommand{name}
}

func RegisterVarintCombiner(name string) {
	streamconv.RegisterCombiner(name, NewVarintCombinerCommand(name))
}
