package converters

import (
	"fmt"
	"io"
	"streamconv"
)

type base64EncoderCommand struct {
	name string
}

func (c *base64EncoderCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *base64EncoderCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewBase64Encoder(), nil
}

func NewBase64EncoderCommand(name string) (command streamconv.ConverterCommand) {
	return &base64EncoderCommand{name}
}

func RegisterBase64Encoder(name string) {
	streamconv.RegisterConverter(name, NewBase64EncoderCommand(name))
}

type base64DecoderCommand struct {
	name string
}

func (c *base64DecoderCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *base64DecoderCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewBase64Decoder(), nil
}

func NewBase64DecoderCommand(name string) (command streamconv.ConverterCommand) {
	return &base64DecoderCommand{name}
}

func RegisterBase64Decoder(name string) {
	streamconv.RegisterConverter(name, NewBase64DecoderCommand(name))
}
