package transformers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/connesc/streamconv"
)

type windowBuffer struct {
	items  []*bytes.Buffer
	length uint
}

func newWindowBuffer(length uint) *windowBuffer {
	items := make([]*bytes.Buffer, length)
	for i := range items {
		items[i] = &bytes.Buffer{}
	}
	return &windowBuffer{items, length}
}

func (b *windowBuffer) get(offset uint) *bytes.Buffer {
	return b.items[offset%b.length]
}

type windowReader struct {
	reader        streamconv.ItemReader
	buffer        *windowBuffer
	bufferIndex   uint
	skippedCount  uint
	reusedCount   uint
	readCount     uint
	bufferedCount uint
	err           error
}

func (r *windowReader) ReadItem() (item io.Reader, err error) {
	if r.err != nil {
		return nil, r.err
	}

	defer func() {
		if err != nil {
			r.err = err
		}
	}()

	for r.skippedCount > 0 {
		var skipped io.Reader
		skipped, err = r.reader.ReadItem()
		if err != nil {
			return
		}
		_, err = io.Copy(ioutil.Discard, skipped)
		if err != nil {
			return
		}
		r.skippedCount--
	}

	if r.reusedCount > 0 {
		item = r.buffer.get(r.bufferIndex)
		r.bufferIndex++
		r.reusedCount--
		return
	}

	if r.readCount == 0 {
		return nil, io.EOF
	}

	item, err = r.reader.ReadItem()
	if err != nil {
		return
	}
	r.readCount--

	if r.bufferedCount > r.readCount {
		item = io.TeeReader(item, r.buffer.get(r.bufferIndex))
		r.bufferIndex++
		r.bufferedCount--
	}
	return
}

type windowExtractor struct {
	reader     streamconv.LookaheadItemReader
	size       uint
	step       uint
	overlap    uint
	hole       uint
	bufferStep uint
	current    *windowReader
	err        error
}

func (e *windowExtractor) ExtractWindow() (window streamconv.ItemReader, err error) {
	if e.err != nil {
		return nil, e.err
	}

	defer func() {
		if err != nil {
			e.err = err
		}
	}()

	_, err = e.reader.Lookahead()
	if err != nil {
		return
	}

	if e.current == nil {
		e.current = &windowReader{
			reader:        e.reader,
			buffer:        newWindowBuffer(e.overlap),
			bufferIndex:   0,
			skippedCount:  0,
			reusedCount:   0,
			readCount:     e.size,
			bufferedCount: e.overlap,
		}
		window = e.current
	} else if e.current.err != io.EOF {
		err = e.current.err
		if err == nil {
			err = fmt.Errorf("the previous window has not been fully read")
		}
	} else if e.current.readCount == 0 {
		e.current.bufferIndex -= e.overlap
		e.current.skippedCount = e.hole
		e.current.reusedCount = e.overlap
		e.current.readCount = e.step
		e.current.bufferedCount = e.bufferStep
		e.current.err = nil
		window = e.current
	} else {
		err = io.EOF
	}
	return
}

type multiWindowReader struct {
	extractor  windowExtractor
	subProgram streamconv.Transformer
	current    streamconv.ItemReader
	err        error
}

func (r *multiWindowReader) ReadItem() (item io.Reader, err error) {
	if r.err != nil {
		return nil, r.err
	}

	defer func() {
		if err != nil {
			r.err = err
		}
	}()

	for {
		if r.current == nil {
			r.current, err = r.extractor.ExtractWindow()
			if err != nil {
				return
			}

			if r.subProgram != nil {
				r.current, err = r.subProgram.Transform(r.current)
				if err != nil {
					if err == io.EOF {
						err = io.ErrUnexpectedEOF
					}
					return
				}
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

type windowTransformer struct {
	subProgram streamconv.Transformer
	size       uint
	step       uint
	partial    bool
}

func (t *windowTransformer) Transform(src streamconv.ItemReader) (dst streamconv.ItemReader, err error) {
	// TODO: use partial

	var overlap, hole, bufferStep uint
	if t.step < t.size {
		overlap = t.size - t.step
		if t.step < overlap {
			bufferStep = t.step
		} else {
			bufferStep = overlap
		}
	} else {
		hole = t.step - t.size
	}

	dst = &multiWindowReader{
		extractor: windowExtractor{
			reader:     streamconv.NewLookaheadItemReader(src),
			size:       t.size,
			step:       t.step,
			overlap:    overlap,
			hole:       hole,
			bufferStep: bufferStep,
		},
		subProgram: t.subProgram,
	}
	return
}

func NewWindowTransformer(subProgram streamconv.Transformer, size uint, step uint, partial bool) streamconv.Transformer {
	return &windowTransformer{
		subProgram: subProgram,
		size:       size,
		step:       step,
		partial:    partial,
	}
}
