package combiners

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type catCombinerCLI struct {
	name string
}

func (c *catCombinerCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *catCombinerCLI) Parse(args []string) (command streamconv.CombinerCommand, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func(out io.Writer) (streamconv.ItemWriter, error) {
		return NewCatCombiner(out), nil
	}
	return
}

func NewCatCombinerCLI(name string) (cli streamconv.CombinerCLI) {
	return &catCombinerCLI{name}
}

func RegisterCatCombiner(name string) {
	streamconv.RegisterCombiner(name, NewCatCombinerCLI(name))
}
