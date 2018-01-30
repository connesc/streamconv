package app

import (
	"io"

	"github.com/connesc/streamconv"
)

type converterTransformer struct {
	command streamconv.ConverterCommand
}

func (t *converterTransformer) Transform(src streamconv.ItemReader) (dst streamconv.ItemReader, err error) {
	converter, err := t.command()
	if err != nil {
		return
	}

	return &converterReader{src, converter}, nil
}

type converterReader struct {
	reader    streamconv.ItemReader
	converter streamconv.Converter
}

func (r *converterReader) ReadItem() (item io.Reader, err error) {
	item, err = r.reader.ReadItem()
	if err != nil {
		return nil, err
	}

	item, err = r.converter.Convert(item)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

type extractorTransformer struct {
	command streamconv.ExtractorCommand
}

func (t *extractorTransformer) Transform(src streamconv.ItemReader) (dst streamconv.ItemReader, err error) {
	return &extractorReader{src, t.command, nil}, nil
}

type extractorReader struct {
	reader    streamconv.ItemReader
	command   streamconv.ExtractorCommand
	extractor streamconv.ItemReader
}

func (r *extractorReader) ReadItem() (item io.Reader, err error) {
	if r.extractor == nil {
		err = io.EOF
	} else {
		item, err = r.extractor.ReadItem()
	}

	for err == io.EOF {
		var src io.Reader
		src, err = r.reader.ReadItem()
		if err != nil {
			return
		}

		r.extractor, err = r.command(src)
		if err != nil {
			return
		}

		item, err = r.extractor.ReadItem()
	}

	return
}

type combinerTransformer struct {
	command streamconv.CombinerCommand
}

func (t *combinerTransformer) Transform(src streamconv.ItemReader) (dst streamconv.ItemReader, err error) {
	return &combinerReader{src, t.command, false}, nil
}

type combinerReader struct {
	reader  streamconv.ItemReader
	command streamconv.CombinerCommand
	done    bool
}

func (r *combinerReader) ReadItem() (item io.Reader, err error) {
	if r.done {
		return nil, io.EOF
	}
	r.done = true

	pr, pw := io.Pipe()
	writer, err := r.command(pw)
	if err != nil {
		return
	}

	go func() {
		err := streamconv.Copy(writer, r.reader)
		pw.CloseWithError(err)
	}()

	return pr, nil
}
