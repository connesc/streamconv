package splitters

import (
	"fmt"
	"io"
	"streamconv"
)

type varintSplitterCommand struct {
	name string
}

func (c *varintSplitterCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *varintSplitterCommand) Parse(args []string, in io.Reader) (splitter streamconv.Splitter, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewVarintSplitter(in), nil
}

func NewVarintSplitterCommand(name string) (command streamconv.SplitterCommand) {
	return &varintSplitterCommand{name}
}

func RegisterVarintSplitter(name string) {
	streamconv.RegisterSplitter(name, NewVarintSplitterCommand(name))
}
