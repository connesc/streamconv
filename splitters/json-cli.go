package splitters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type jsonSplitterCommand struct {
	name string
}

func (c *jsonSplitterCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *jsonSplitterCommand) Parse(args []string, in io.Reader) (splitter streamconv.Splitter, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewJSONSplitter(in), nil
}

func NewJSONSplitterCommand(name string) (command streamconv.SplitterCommand) {
	return &jsonSplitterCommand{name}
}

func RegisterJSONSplitter(name string) {
	streamconv.RegisterSplitter(name, NewJSONSplitterCommand(name))
}
