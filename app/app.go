package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/connesc/streamconv"
	"github.com/connesc/streamconv/combiners"
	"github.com/connesc/streamconv/extractors"
	"github.com/connesc/streamconv/parser"
)

type App interface {
	Run(dst io.Writer, src io.Reader) error
}

type streamconvApp struct {
	extractor    streamconv.ExtractorCommand
	transformers []streamconv.Transformer
	combiner     streamconv.CombinerCommand
}

var defaultExtractor = func(in io.Reader) (streamconv.ItemReader, error) {
	return extractors.NewSingleExtractor(in), nil
}

var defaultCombiner = func(out io.Writer) (streamconv.ItemWriter, error) {
	return combiners.NewCatCombiner(out), nil
}

func Parse(program string) (app App, err error) {
	parsedProgram, err := parser.Parse(strings.NewReader(program))
	if err != nil {
		return
	}
	return New(parsedProgram)
}

func New(program parser.Program) (app App, err error) {
	var extractor streamconv.ExtractorCommand
	var transformers []streamconv.Transformer
	var combiner streamconv.CombinerCommand

	for _, commandSource := range program {
		command, err := streamconv.ParseCommand(commandSource)
		if err != nil {
			return nil, err
		}

		if combiner != nil {
			transformers = append(transformers, &combinerTransformer{combiner})
			combiner = nil
		}

		if extractor == nil {
			if command, ok := command.(streamconv.ExtractorCommand); ok {
				extractor = command
				continue
			} else {
				extractor = defaultExtractor
			}
		}

		switch command := command.(type) {
		case streamconv.ExtractorCommand:
			transformers = append(transformers, &extractorTransformer{command})
		case streamconv.ConverterCommand:
			transformers = append(transformers, &converterTransformer{command})
		case streamconv.CombinerCommand:
			combiner = command
		default:
			return nil, fmt.Errorf("unknown command type: %T", command)
		}
	}

	if extractor == nil {
		extractor = defaultExtractor
	}

	if combiner == nil {
		combiner = defaultCombiner
	}

	return &streamconvApp{extractor, transformers, combiner}, nil
}

func (app streamconvApp) Run(dst io.Writer, src io.Reader) (err error) {
	reader, err := app.extractor(src)
	if err != nil {
		return
	}

	for _, transformer := range app.transformers {
		reader, err = transformer.Transform(reader)
		if err != nil {
			return err
		}
	}

	writer, err := app.combiner(dst)
	if err != nil {
		return
	}

	return streamconv.Copy(writer, reader)
}
