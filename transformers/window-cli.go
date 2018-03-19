package transformers

import (
	"fmt"
	"io"
	"strconv"

	"github.com/connesc/streamconv"

	"github.com/spf13/pflag"
)

type windowTransformerCLI struct {
	name string
}

type windowTransformerOptions struct {
	skipPartial bool
}

func (c *windowTransformerCLI) newFlagSet() (flags *pflag.FlagSet, options *windowTransformerOptions) {
	options = &windowTransformerOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.BoolVar(&options.skipPartial, "skip-partial", false, "skip the eventual partial group at the end")
	return
}

func (c *windowTransformerCLI) PrintUsage(output io.Writer) (err error) {
	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *windowTransformerCLI) Parse(args []string, subProgram streamconv.TransformerCommand) (command streamconv.TransformerCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	args = flags.Args()
	switch {
	case len(args) == 0:
		return nil, fmt.Errorf("at least one argument is required")
	case len(args) > 2:
		return nil, fmt.Errorf("too many arguments (expected up to 2, got %v)", len(args))
	}

	size, err := strconv.ParseUint(args[0], 10, 31)
	if err != nil {
		return
	}

	step := size
	if len(args) == 2 {
		step, err = strconv.ParseUint(args[1], 10, 31)
	}
	if err != nil {
		return
	}

	var subTransformer streamconv.Transformer
	if subProgram != nil {
		subTransformer, err = subProgram()
		if err != nil {
			return
		}
	}

	command = func() (streamconv.Transformer, error) {
		return NewWindowTransformer(subTransformer, uint(size), uint(step), options.skipPartial), nil
	}
	return
}

func NewWindowTransformerCLI(name string) (cli streamconv.TransformerCLI) {
	return &windowTransformerCLI{name}
}

func RegisterWindowTransformer(name string) {
	streamconv.RegisterTransformer(name, NewWindowTransformerCLI(name))
}
