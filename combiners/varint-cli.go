package combiners

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type varintCombinerCLI struct {
	name string
}

func (c *varintCombinerCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *varintCombinerCLI) Parse(args []string) (command streamconv.CombinerCommand, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func(out io.Writer) (streamconv.ItemWriter, error) {
		return NewVarintCombiner(out), nil
	}
	return
}

func NewVarintCombinerCLI(name string) (cli streamconv.CombinerCLI) {
	return &varintCombinerCLI{name}
}

func RegisterVarintCombiner(name string) {
	streamconv.RegisterCombiner(name, NewVarintCombinerCLI(name))
}
