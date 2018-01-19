package converters

import (
	"fmt"
	"io"
	"streamconv"
)

type protobufToJSONCommand struct {
	name string
}

func (b *protobufToJSONCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (b *protobufToJSONCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	return NewProtobufToJSON(args[0], args[1]), nil
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

func (b *protobufFromJSONCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (b *protobufFromJSONCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	return NewProtobufFromJSON(args[0], args[1]), nil
}

func NewProtobufFromJSONCommand(name string) (command streamconv.ConverterCommand) {
	return &protobufFromJSONCommand{name}
}

func RegisterProtobufFromJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufFromJSONCommand(name))
}
