package streamconv

import (
	"io"
)

type CLI interface {
	PrintUsage(output io.Writer) (err error)
}

type ExtractorCLI interface {
	CLI
	Parse(args []string) (command ExtractorCommand, err error)
}

type ConverterCLI interface {
	CLI
	Parse(args []string) (command ConverterCommand, err error)
}

type CombinerCLI interface {
	CLI
	Parse(args []string) (command CombinerCommand, err error)
}

type ExtractorCommand interface {
	Run(in io.Reader) (extractor ItemReader, err error)
}

type ConverterCommand interface {
	Run() (converter Converter, err error)
}

type CombinerCommand interface {
	Run(out io.Writer) (combiner ItemWriter, err error)
}
