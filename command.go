package streamconv

import (
	"io"
)

type Command interface {
	PrintUsage(output io.Writer) (err error)
}

type SplitterCommand interface {
	Command
	Parse(args []string, in io.Reader) (splitter Splitter, err error)
}

type ConverterCommand interface {
	Command
	Parse(args []string) (converter Converter, err error)
}

type JoinerCommand interface {
	Command
	Parse(args []string, out io.Writer) (joiner Joiner, err error)
}
