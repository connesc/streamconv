package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type protobufToJSONCommand struct {
	protoFile   string
	messageName string
}

func (c *protobufToJSONCommand) Run() (converter streamconv.Converter, err error) {
	return NewProtobufToJSON(c.protoFile, c.messageName)
}

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

	return &protobufToJSONCommand{args[0], args[1]}, nil
}

func NewProtobufToJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufToJSONCLI{name}
}

func RegisterProtobufToJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufToJSONCLI(name))
}

type protobufFromJSONCommand struct {
	protoFile   string
	messageName string
}

func (c *protobufFromJSONCommand) Run() (converter streamconv.Converter, err error) {
	return NewProtobufFromJSON(c.protoFile, c.messageName)
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

	return &protobufFromJSONCommand{args[0], args[1]}, nil
}

func NewProtobufFromJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufFromJSONCLI{name}
}

func RegisterProtobufFromJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufFromJSONCLI(name))
}
