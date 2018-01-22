package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type streamConverterCommand struct {
	name string
}

func (c *streamConverterCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *streamConverterCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("invalid number of arguments (exepcted 1, got %v)", len(args))
	}

	return NewStreamConverter(args[0])
}

func NewStreamConverterCommand(name string) (command streamconv.ConverterCommand) {
	return &streamConverterCommand{name}
}

func RegisterStreamConverter(name string) {
	streamconv.RegisterConverter(name, NewStreamConverterCommand(name))
}
