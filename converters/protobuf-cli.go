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

type protobufToTextCLI struct {
	name string
}

func (c *protobufToTextCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufToTextCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufToText(args[0], args[1])
	}
	return
}

func NewProtobufToTextCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufToTextCLI{name}
}

func RegisterProtobufToText(name string) {
	streamconv.RegisterConverter(name, NewProtobufToTextCLI(name))
}

type protobufFromTextCLI struct {
	name string
}

func (c *protobufFromTextCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufFromTextCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufFromText(args[0], args[1])
	}
	return
}

func NewProtobufFromTextCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufFromTextCLI{name}
}

func RegisterProtobufFromText(name string) {
	streamconv.RegisterConverter(name, NewProtobufFromTextCLI(name))
}

type protobufJSONToTextCLI struct {
	name string
}

func (c *protobufJSONToTextCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufJSONToTextCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufJSONToText(args[0], args[1])
	}
	return
}

func NewProtobufJSONToTextCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufJSONToTextCLI{name}
}

func RegisterProtobufJSONToText(name string) {
	streamconv.RegisterConverter(name, NewProtobufJSONToTextCLI(name))
}

type protobufTextToJSONCLI struct {
	name string
}

func (c *protobufTextToJSONCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *protobufTextToJSONCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufTextToJSON(args[0], args[1])
	}
	return
}

func NewProtobufTextToJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufTextToJSONCLI{name}
}

func RegisterProtobufTextToJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufTextToJSONCLI(name))
}
