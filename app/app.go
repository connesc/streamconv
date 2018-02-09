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
	Transformer() streamconv.Transformer
}

type streamconvApp struct {
	extractor   streamconv.ExtractorCommand
	transformer streamconv.Transformer
	combiner    streamconv.CombinerCommand
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
		var subProgram streamconv.TransformerCommand
		if len(commandSource.SubProgram) > 0 {
			subProgram = func() (streamconv.Transformer, error) {
				app, err := New(commandSource.SubProgram)
				if err != nil {
					return nil, err
				}
				return app.Transformer(), nil
			}
		}

		command, err := streamconv.ParseCommand(commandSource.Words, subProgram)
		if err != nil {
			return nil, err
		}

		if combiner != nil {
			transformers = append(transformers, streamconv.NewCombinerTransformer(combiner))
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
			transformers = append(transformers, streamconv.NewExtractorTransformer(command))
		case streamconv.ConverterCommand:
			transformers = append(transformers, streamconv.NewConverterTransformer(command))
		case streamconv.TransformerCommand:
			transformer, err := command()
			if err != nil {
				return nil, err
			}
			transformers = append(transformers, transformer)
		case streamconv.CombinerCommand:
			combiner = command
		default:
			return nil, fmt.Errorf("unknown command type: %T", command)
		}
	}

	return &streamconvApp{extractor, streamconv.Compose(transformers...), combiner}, nil
}

func (app *streamconvApp) Run(dst io.Writer, src io.Reader) (err error) {
	extractor := app.extractor
	if extractor == nil {
		extractor = defaultExtractor
	}
	extractor = streamconv.Transform(extractor, app.transformer)

	combiner := app.combiner
	if combiner == nil {
		combiner = defaultCombiner
	}

	reader, err := extractor(src)
	if err != nil {
		return
	}

	writer, err := combiner(dst)
	if err != nil {
		return
	}

	return streamconv.Copy(writer, reader)
}

func (app *streamconvApp) Transformer() streamconv.Transformer {
	transformers := make([]streamconv.Transformer, 0, 3)
	if app.extractor != nil {
		transformers = append(transformers, streamconv.NewExtractorTransformer(app.extractor))
	}
	transformers = append(transformers, app.transformer)
	if app.combiner != nil {
		transformers = append(transformers, streamconv.NewCombinerTransformer(app.combiner))
	}

	return streamconv.Compose(transformers...)
}
