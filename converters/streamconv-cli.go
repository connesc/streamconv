package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

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

	command = func() (streamconv.Converter, error) {
		return NewStreamConverter(args[0])
	}
	return
}

func NewStreamConverterCLI(name string) (cli streamconv.ConverterCLI) {
	return &streamConverterCLI{name}
}

func RegisterStreamConverter(name string) {
	streamconv.RegisterConverter(name, NewStreamConverterCLI(name))
}
