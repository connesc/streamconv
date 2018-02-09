package streamconv

import (
	"io"
)

func Copy(dst ItemWriter, src ItemReader) (err error) {
	for {
		item, err := src.ReadItem()
		if err != nil {
			return err
		}

		err = dst.WriteItem(item)
		if err != nil {
			return err
		}
	}
}

func Compose(transformers ...Transformer) (composed Transformer) {
	return &composedTransformer{transformers}
}

type composedTransformer struct {
	transformers []Transformer
}

func (t *composedTransformer) Transform(src ItemReader) (dst ItemReader, err error) {
	dst = src
	for _, transformer := range t.transformers {
		dst, err = transformer.Transform(dst)
		if err != nil {
			return
		}
	}
	return
}

func NewConverterTransformer(converter ConverterCommand) Transformer {
	return &converterTransformer{converter}
}

type converterTransformer struct {
	command ConverterCommand
}

func (t *converterTransformer) Transform(src ItemReader) (dst ItemReader, err error) {
	converter, err := t.command()
	if err != nil {
		return
	}

	return &converterReader{src, converter}, nil
}

type converterReader struct {
	reader    ItemReader
	converter Converter
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

func NewExtractorTransformer(extractor ExtractorCommand) Transformer {
	return &extractorTransformer{extractor}
}

type extractorTransformer struct {
	command ExtractorCommand
}

func (t *extractorTransformer) Transform(src ItemReader) (dst ItemReader, err error) {
	return &extractorReader{src, t.command, nil}, nil
}

type extractorReader struct {
	reader    ItemReader
	command   ExtractorCommand
	extractor ItemReader
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

func NewCombinerTransformer(combiner CombinerCommand) Transformer {
	return &combinerTransformer{combiner}
}

type combinerTransformer struct {
	command CombinerCommand
}

func (t *combinerTransformer) Transform(src ItemReader) (dst ItemReader, err error) {
	return &combinerReader{src, t.command, false}, nil
}

type combinerReader struct {
	reader  ItemReader
	command CombinerCommand
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
		err := Copy(writer, r.reader)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func Transform(command ExtractorCommand, transformer Transformer) (transformed ExtractorCommand) {
	return func(in io.Reader) (extractor ItemReader, err error) {
		extractor, err = command(in)
		if err != nil {
			return
		}
		return transformer.Transform(extractor)
	}
}

func TransformItems(transformer Transformer) (transformed Transformer) {
	return &itemTransformer{
		transformer: transformer,
	}
}

type itemTransformer struct {
	transformer Transformer
}

func (t *itemTransformer) Transform(src ItemReader) (dst ItemReader, err error) {
	dst = &transformedItemReader{
		reader:      src,
		transformer: t.transformer,
	}
	return
}

type transformedItemReader struct {
	reader      ItemReader
	transformer Transformer
	current     ItemReader
	err         error
}

func (r *transformedItemReader) ReadItem() (item io.Reader, err error) {
	if r.err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			r.err = err
		}
	}()

	for {
		if r.current == nil {
			var src io.Reader
			src, err = r.reader.ReadItem()
			if err != nil {
				return
			}

			r.current, err = r.transformer.Transform(&singleItemReader{item: src})
			if err != nil {
				return
			}
		}

		item, err = r.current.ReadItem()
		if err != io.EOF {
			return
		}
		err = nil
		r.current = nil
	}
}

type singleItemReader struct {
	item io.Reader
	done bool
}

func (r *singleItemReader) ReadItem() (item io.Reader, err error) {
	if r.done {
		return nil, io.EOF
	}
	r.done = true
	return r.item, nil
}

type LookaheadItemReader interface {
	ItemReader
	Lookahead() (head io.Reader, err error)
}

func NewLookaheadItemReader(reader ItemReader) LookaheadItemReader {
	return &lookaheadItemReader{
		reader: reader,
	}
}

type lookaheadItemReader struct {
	reader ItemReader
	head   io.Reader
	err    error
}

func (r *lookaheadItemReader) ReadItem() (item io.Reader, err error) {
	switch {
	case r.err != nil:
		err = r.err
	case r.head != nil:
		item = r.head
		r.head = nil
	default:
		item, err = r.reader.ReadItem()
	}
	return
}

func (r *lookaheadItemReader) Lookahead() (head io.Reader, err error) {
	switch {
	case r.err != nil:
		err = r.err
	case r.head != nil:
		head = r.head
	default:
		head, err = r.reader.ReadItem()
		r.head = head
	}
	return
}
