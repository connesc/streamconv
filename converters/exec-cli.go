package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type executorCommand struct {
	name string
	args []string
}

func (c *executorCommand) Run() (converter streamconv.Converter, err error) {
	return NewExecutor(c.name, c.args), nil
}

type executorCLI struct {
	name string
}

func (c *executorCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *executorCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("at least one argument is expected")
	}

	return &executorCommand{args[0], args[1:]}, nil
}

func NewExecutorCLI(name string) (cli streamconv.ConverterCLI) {
	return &executorCLI{name}
}

func RegisterExecutor(name string) {
	streamconv.RegisterConverter(name, NewExecutorCLI(name))
}
