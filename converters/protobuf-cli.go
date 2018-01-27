package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type protobufToJSONCLI struct {
	name string
}

func (c *protobufToJSONCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufToJSONCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufToJSON(args[0], args[1])
	}
	return
}

func NewProtobufToJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufToJSONCLI{name}
}

func RegisterProtobufToJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufToJSONCLI(name))
}

type protobufFromJSONCLI struct {
	name string
}

func (c *protobufFromJSONCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufFromJSONCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufFromJSON(args[0], args[1])
	}
	return
}

func NewProtobufFromJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufFromJSONCLI{name}
}

func RegisterProtobufFromJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufFromJSONCLI(name))
}
