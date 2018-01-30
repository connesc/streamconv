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

type TransformerCLI interface {
	CLI
	Parse(args []string) (command TransformerCommand, err error)
}

type CombinerCLI interface {
	CLI
	Parse(args []string) (command CombinerCommand, err error)
}

type ExtractorCommand func(in io.Reader) (extractor ItemReader, err error)

type ConverterCommand func() (converter Converter, err error)

type TransformerCommand func() (transformer Transformer, err error)

type CombinerCommand func(out io.Writer) (combiner ItemWriter, err error)
