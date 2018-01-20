package converters

import (
	"fmt"
	"io"
	"streamconv"
)

type executorCommand struct {
	name string
}

func (b *executorCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (b *executorCommand) Parse(args []string) (converter streamconv.Converter, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("at least one argument is expected")
	}

	return NewExecutor(args[0], args[1:]), nil
}

func NewExecutorCommand(name string) (command streamconv.ConverterCommand) {
	return &executorCommand{name}
}

func RegisterExecutorCommand(name string) {
	streamconv.RegisterConverter(name, NewExecutorCommand(name))
}
