package app

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
	"github.com/connesc/streamconv/combiners"
	"github.com/connesc/streamconv/extractors"
)

type App interface {
	Run(dst io.Writer, src io.Reader) error
}

type streamConverter interface {
	Convert(src streamconv.ItemReader) (dst streamconv.ItemReader, err error)
}

type streamconvApp struct {
	extractor  streamconv.ExtractorCommand
	converters []streamConverter
	combiner   streamconv.CombinerCommand
}

var defaultExtractor = func(in io.Reader) (streamconv.ItemReader, error) {
	return extractors.NewSingleExtractor(in), nil
}

var defaultCombiner = func(out io.Writer) (streamconv.ItemWriter, error) {
	return combiners.NewJoinCombiner(out, ""), nil
}

func New(program string) (app App, err error) {
	commands, err := parse(program)

	var extractor streamconv.ExtractorCommand
	var converters []streamConverter
	var combiner streamconv.CombinerCommand

	for _, tokens := range commands {
		command, err := streamconv.ParseCommand(tokens)
		if err != nil {
			return nil, err
		}

		if combiner != nil {
			converters = append(converters, &combinerConverter{combiner})
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
			converters = append(converters, &extractorConverter{command})
		case streamconv.ConverterCommand:
			converters = append(converters, &regularConverter{command})
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

	return &streamconvApp{extractor, converters, combiner}, nil
}

func (app streamconvApp) Run(dst io.Writer, src io.Reader) (err error) {
	reader, err := app.extractor(src)
	if err != nil {
		return
	}

	for _, converter := range app.converters {
		reader, err = converter.Convert(reader)
	}

	writer, err := app.combiner(dst)
	if err != nil {
		return
	}

	return streamconv.Copy(writer, reader)
}
