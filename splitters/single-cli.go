package splitters

import (
	"fmt"
	"io"
	"streamconv"
)

type singleSplitterCommand struct {
	name string
}

func (c *singleSplitterCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *singleSplitterCommand) Parse(args []string, in io.Reader) (splitter streamconv.Splitter, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewSingleSplitter(in), nil
}

func NewSingleSplitterCommand(name string) (command streamconv.SplitterCommand) {
	return &singleSplitterCommand{name}
}

func RegisterSingleSplitter(name string) {
	streamconv.RegisterSplitter(name, NewSingleSplitterCommand(name))
}
