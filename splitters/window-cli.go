package splitters

import (
	"fmt"
	"io"
	"strconv"
	"streamconv"

	"github.com/spf13/pflag"
)

type windowSplitterCommand struct {
	name string
}

type windowSplitterOptions struct {
	partial bool
}

func (c *windowSplitterCommand) newFlagSet() (flags *pflag.FlagSet, options *windowSplitterOptions) {
	options = &windowSplitterOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.BoolVar(&options.partial, "partial", false, "include partial groups")
	return
}

func (c *windowSplitterCommand) PrintUsage(output io.Writer) (err error) {
	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *windowSplitterCommand) Parse(args []string, in io.Reader) (splitter streamconv.Splitter, err error) {
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

	return NewWindowSplitter(in, int(size), int(step), options.partial), nil
}

func NewWindowSplitterCommand(name string) (command streamconv.SplitterCommand) {
	return &windowSplitterCommand{name}
}

func RegisterWindowSplitter(name string) {
	streamconv.RegisterSplitter(name, NewWindowSplitterCommand(name))
}
