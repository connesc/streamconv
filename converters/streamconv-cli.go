package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type streamConverterCommand struct {
	program string
}

func (c *streamConverterCommand) Run() (converter streamconv.Converter, err error) {
	return NewStreamConverter(c.program)
}

type streamConverterCLI struct {
	name string
}

func (c *streamConverterCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *streamConverterCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments (exepcted 1, got %v)", len(args))
	}

	return &streamConverterCommand{args[0]}, nil
}

func NewStreamConverterCLI(name string) (cli streamconv.ConverterCLI) {
	return &streamConverterCLI{name}
}

func RegisterStreamConverter(name string) {
	streamconv.RegisterConverter(name, NewStreamConverterCLI(name))
}
