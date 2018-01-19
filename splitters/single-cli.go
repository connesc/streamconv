package splitters

import (
	"fmt"
	"io"
	"streamconv"
)

type singleSplitterCommand struct {
	name string
}

func (b *singleSplitterCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (b *singleSplitterCommand) Parse(args []string, in io.Reader) (splitter streamconv.Splitter, err error) {
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
