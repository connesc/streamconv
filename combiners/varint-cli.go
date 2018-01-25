package combiners

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type varintCombinerCommand struct{}

func (c *varintCombinerCommand) Run(out io.Writer) (combiner streamconv.ItemWriter, err error) {
	return NewVarintCombiner(out), nil
}

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

	return &varintCombinerCommand{}, nil
}

func NewVarintCombinerCLI(name string) (cli streamconv.CombinerCLI) {
	return &varintCombinerCLI{name}
}

func RegisterVarintCombiner(name string) {
	streamconv.RegisterCombiner(name, NewVarintCombinerCLI(name))
}
