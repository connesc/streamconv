package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type protobufToJSONCommand struct {
	name string
}

func (c *protobufToJSONCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufToJSONCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	return NewProtobufToJSON(args[0], args[1])
}

func NewProtobufToJSONCommand(name string) (command streamconv.ConverterCommand) {
	return &protobufToJSONCommand{name}
}

func RegisterProtobufToJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufToJSONCommand(name))
}

type protobufFromJSONCommand struct {
	name string
}

func (c *protobufFromJSONCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufFromJSONCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	return NewProtobufFromJSON(args[0], args[1])
}

func NewProtobufFromJSONCommand(name string) (command streamconv.ConverterCommand) {
	return &protobufFromJSONCommand{name}
}

func RegisterProtobufFromJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufFromJSONCommand(name))
}
