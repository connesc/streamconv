package combiners

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type joinCombinerCommand struct {
	name string
}

func (c *joinCombinerCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *joinCombinerCommand) Parse(args []string, out io.Writer) (combiner streamconv.ItemWriter, err error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments (expected up to 1, got %v)", len(args))
	}

	delimiter := "\n"
	if len(args) == 1 {
		delimiter = args[0]
	}

	return NewJoinCombiner(out, delimiter), nil
}

func NewJoinCombinerCommand(name string) (command streamconv.CombinerCommand) {
	return &joinCombinerCommand{name}
}

func RegisterJoinCombiner(name string) {
	streamconv.RegisterCombiner(name, NewJoinCombinerCommand(name))
}
