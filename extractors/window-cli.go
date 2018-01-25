package extractors

import (
	"fmt"
	"io"
	"strconv"

	"github.com/connesc/streamconv"

	"github.com/spf13/pflag"
)

type windowExtractorCommand struct {
	step    int
	size    int
	partial bool
}

func (c *windowExtractorCommand) Run(in io.Reader) (extractor streamconv.ItemReader, err error) {
	return NewWindowExtractor(in, c.step, c.size, c.partial), nil
}

type windowExtractorCLI struct {
	name string
}

type windowExtractorOptions struct {
	partial bool
}

func (c *windowExtractorCLI) newFlagSet() (flags *pflag.FlagSet, options *windowExtractorOptions) {
	options = &windowExtractorOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.BoolVar(&options.partial, "partial", false, "include partial groups")
	return
}

func (c *windowExtractorCLI) PrintUsage(output io.Writer) (err error) {
	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *windowExtractorCLI) Parse(args []string) (command streamconv.ExtractorCommand, err error) {
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

	return &windowExtractorCommand{int(size), int(step), options.partial}, nil
}

func NewWindowExtractorCLI(name string) (cli streamconv.ExtractorCLI) {
	return &windowExtractorCLI{name}
}

func RegisterWindowExtractor(name string) {
	streamconv.RegisterExtractor(name, NewWindowExtractorCLI(name))
}
