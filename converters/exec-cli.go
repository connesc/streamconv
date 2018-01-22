package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type executorCommand struct {
	name string
}

func (c *executorCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *executorCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("at least one argument is expected")
	}

	return NewExecutor(args[0], args[1:]), nil
}

func NewExecutorCommand(name string) (command streamconv.ConverterCommand) {
	return &executorCommand{name}
}

func RegisterExecutor(name string) {
	streamconv.RegisterConverter(name, NewExecutorCommand(name))
}
