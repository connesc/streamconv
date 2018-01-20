package splitters

import (
	"fmt"
	"io"
	"regexp"
	"streamconv"
)

type simpleSplitterCommand struct {
	name string
}

func (b *simpleSplitterCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (b *simpleSplitterCommand) Parse(args []string, in io.Reader) (splitter streamconv.Splitter, err error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments (expected up to 1, got %v)", len(args))
	}

	rawDelimiter := "\\n"
	if len(args) == 1 {
		rawDelimiter = args[0]
	}

	delimiter, err := regexp.Compile(rawDelimiter)
	if err != nil {
		return
	}

	return NewSimpleSplitter(in, delimiter), nil
}

func NewSimpleSplitterCommand(name string) (command streamconv.SplitterCommand) {
	return &simpleSplitterCommand{name}
}

func RegisterSimpleSplitter(name string) {
	streamconv.RegisterSplitter(name, NewSimpleSplitterCommand(name))
}