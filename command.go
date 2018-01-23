package streamconv

import (
	"io"
)

type Command interface {
	PrintUsage(output io.Writer) (err error)
}

type ExtractorCommand interface {
	Command
	Parse(args []string, in io.Reader) (extractor ItemReader, err error)
}

type ConverterCommand interface {
	Command
	Parse(args []string) (converter Converter, err error)
}

type CombinerCommand interface {
	Command
	Parse(args []string, out io.Writer) (combiner ItemWriter, err error)
}
