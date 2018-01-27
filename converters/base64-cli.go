package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type base64EncoderCLI struct {
	name string
}

func (c *base64EncoderCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *base64EncoderCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewBase64Encoder(), nil
	}
	return
}

func NewBase64EncoderCLI(name string) (cli streamconv.ConverterCLI) {
	return &base64EncoderCLI{name}
}

func RegisterBase64Encoder(name string) {
	streamconv.RegisterConverter(name, NewBase64EncoderCLI(name))
}

type base64DecoderCLI struct {
	name string
}

func (c *base64DecoderCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *base64DecoderCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewBase64Decoder(), nil
	}
	return
}

func NewBase64DecoderCLI(name string) (cli streamconv.ConverterCLI) {
	return &base64DecoderCLI{name}
}

func RegisterBase64Decoder(name string) {
	streamconv.RegisterConverter(name, NewBase64DecoderCLI(name))
}
