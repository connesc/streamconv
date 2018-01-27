package combiners

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type joinCombinerCLI struct {
	name string
}

func (c *joinCombinerCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *joinCombinerCLI) Parse(args []string) (command streamconv.CombinerCommand, err error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments (expected up to 1, got %v)", len(args))
	}

	delimiter := "\n"
	if len(args) == 1 {
		delimiter = args[0]
	}

	command = func(out io.Writer) (streamconv.ItemWriter, error) {
		return NewJoinCombiner(out, delimiter), nil
	}
	return
}

func NewJoinCombinerCLI(name string) (cli streamconv.CombinerCLI) {
	return &joinCombinerCLI{name}
}

func RegisterJoinCombiner(name string) {
	streamconv.RegisterCombiner(name, NewJoinCombinerCLI(name))
}
