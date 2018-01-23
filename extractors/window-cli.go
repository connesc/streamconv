package extractors

import (
	"fmt"
	"io"
	"strconv"

	"github.com/connesc/streamconv"

	"github.com/spf13/pflag"
)

type windowExtractorCommand struct {
	name string
}

type windowExtractorOptions struct {
	partial bool
}

func (c *windowExtractorCommand) newFlagSet() (flags *pflag.FlagSet, options *windowExtractorOptions) {
	options = &windowExtractorOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.BoolVar(&options.partial, "partial", false, "include partial groups")
	return
}

func (c *windowExtractorCommand) PrintUsage(output io.Writer) (err error) {
	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *windowExtractorCommand) Parse(args []string, in io.Reader) (extractor streamconv.ItemReader, err error) {
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

	return NewWindowExtractor(in, int(size), int(step), options.partial), nil
}

func NewWindowExtractorCommand(name string) (command streamconv.ExtractorCommand) {
	return &windowExtractorCommand{name}
}

func RegisterWindowExtractor(name string) {
	streamconv.RegisterExtractor(name, NewWindowExtractorCommand(name))
}
