package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
	"github.com/spf13/pflag"
)

type jsonIndenterCLI struct {
	name string
}

type jsonIndenterOptions struct {
	prefix string
}

func (c *jsonIndenterCLI) newFlagSet() (flags *pflag.FlagSet, options *jsonIndenterOptions) {
	options = &jsonIndenterOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringVarP(&options.prefix, "prefix", "p", "", "a string to prepend to each line except the first")
	return
}

func (c *jsonIndenterCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	if err != nil {
		return
	}

	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *jsonIndenterCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	indent := "  "

	args = flags.Args()
	switch {
	case len(args) == 1:
		indent = args[0]
	case len(args) != 0:
		return nil, fmt.Errorf("too many arguments (expected up to 1, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewJSONIndenter(options.prefix, indent), nil
	}
	return
}

func NewJSONIndenterCLI(name string) (cli streamconv.ConverterCLI) {
	return &jsonIndenterCLI{name}
}

func RegisterJSONIndenter(name string) {
	streamconv.RegisterConverter(name, NewJSONIndenterCLI(name))
}

type jsonCompactorCLI struct {
	name string
}

func (c *jsonCompactorCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *jsonCompactorCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) != 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewJSONCompactor(), nil
	}
	return
}

func NewJSONCompactorCLI(name string) (cli streamconv.ConverterCLI) {
	return &jsonCompactorCLI{name}
}

func RegisterJSONCompactor(name string) {
	streamconv.RegisterConverter(name, NewJSONCompactorCLI(name))
}
