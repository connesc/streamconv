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

func skip(reader streamconv.ItemReader, count uint) (err error) {
	for count > 0 {
		var skipped io.Reader
		skipped, err = reader.ReadItem()
		if err != nil {
			return
		}
		_, err = io.Copy(ioutil.Discard, skipped)
		if err != nil {
			return
		}
		count--
	}
	return
}

type windowExtractor interface {
	ExtractWindow() (window streamconv.ItemReader, err error)
}

type partialWindowReader struct {
	reader        streamconv.ItemReader
	buffer        *windowBuffer
	bufferIndex   uint
	reusedCount   uint
	readCount     uint
	bufferedCount uint
	err           error
}

func (r *partialWindowReader) ReadItem() (item io.Reader, err error) {
	if r.err != nil {
		return nil, r.err
	}

	defer func() {
		if err != nil {
			r.err = err
		}
	}()

	if r.reusedCount > 0 {
		item = bytes.NewReader(r.buffer.get(r.bufferIndex).Bytes())
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
		buffer := r.buffer.get(r.bufferIndex)
		buffer.Reset()
		item = io.TeeReader(item, buffer)
		r.bufferIndex++
		r.bufferedCount--
	}
	return
}

type partialWindowExtractor struct {
	reader     streamconv.LookaheadItemReader
	size       uint
	step       uint
	overlap    uint
	hole       uint
	bufferStep uint
	current    *partialWindowReader
	err        error
}

func (e *partialWindowExtractor) ExtractWindow() (window streamconv.ItemReader, err error) {
	if e.err != nil {
		return nil, e.err
	}

	defer func() {
		if err != nil {
			e.err = err
		}
	}()

	if e.current == nil {
		_, err = e.reader.Lookahead()
		if err != nil {
			return
		}
		e.current = &partialWindowReader{
			reader:        e.reader,
			buffer:        newWindowBuffer(e.overlap),
			bufferIndex:   0,
			reusedCount:   0,
			readCount:     e.size,
			bufferedCount: e.overlap,
		}
		window = e.current
	} else if err = e.current.err; err != io.EOF {
		if err == nil {
			err = fmt.Errorf("the previous window has not been fully read")
		}
	} else {
		err = skip(e.reader, e.hole)
		if err != nil {
			return
		}
		_, err = e.reader.Lookahead()
		if err != nil {
			return
		}
		e.current.bufferIndex -= e.overlap
		e.current.reusedCount = e.overlap
		e.current.readCount = e.size - e.overlap
		e.current.bufferedCount = e.bufferStep
		e.current.err = nil
		window = e.current
	}
	return
}

type completeWindowReader struct {
	reader      streamconv.ItemReader
	buffer      *windowBuffer
	bufferIndex uint
	forwardLast bool
	lastItem    io.Reader
	index       uint
	err         error
}

func (r *completeWindowReader) fill() (err error) {
	for i := uint(r.index); i < r.buffer.length; i++ {
		var item io.Reader
		item, err = r.reader.ReadItem()
		if err != nil {
			return
		}
		buffer := r.buffer.get(r.bufferIndex + i)
		buffer.Reset()
		_, err = io.Copy(buffer, item)
		if err != nil {
			return
		}
	}
	r.index = 0

	if r.forwardLast {
		r.lastItem, err = r.reader.ReadItem()
		if err != nil {
			return
		}
	}

	return
}

func (r *completeWindowReader) ReadItem() (item io.Reader, err error) {
	if r.err != nil {
		return nil, r.err
	}

	defer func() {
		if err != nil {
			r.err = err
		}
	}()

	if r.index == r.buffer.length {
		item = r.lastItem
		r.lastItem = nil
		r.err = io.EOF
		return
	}

	item = bytes.NewReader(r.buffer.get(r.bufferIndex).Bytes())
	r.bufferIndex++
	r.index++

	if (r.index == r.buffer.length) && (r.lastItem == nil) {
		r.err = io.EOF
	}
	return
}

type completeWindowExtractor struct {
	reader  streamconv.LookaheadItemReader
	size    uint
	step    uint
	overlap uint
	hole    uint
	current *completeWindowReader
	err     error
}

func (e *completeWindowExtractor) ExtractWindow() (window streamconv.ItemReader, err error) {
	if e.err != nil {
		return nil, e.err
	}

	defer func() {
		if err != nil {
			e.err = err
		}
	}()

	if e.current == nil {
		forwardLast := e.overlap == 0
		bufferLength := e.size
		if forwardLast {
			bufferLength--
		}
		e.current = &completeWindowReader{
			reader:      e.reader,
			buffer:      newWindowBuffer(bufferLength),
			bufferIndex: 0,
			forwardLast: forwardLast,
			index:       0,
		}
	} else if err = e.current.err; err != io.EOF {
		if err == nil {
			err = fmt.Errorf("the previous window has not been fully read")
		}
		return
	} else {
		err = skip(e.reader, e.hole)
		if err != nil {
			return
		}
		e.current.bufferIndex -= e.overlap
		e.current.index = e.overlap
		e.current.err = nil
	}

	err = e.current.fill()
	if err == nil {
		window = e.current
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
	subProgram  streamconv.Transformer
	size        uint
	step        uint
	skipPartial bool
}

func (t *windowTransformer) Transform(src streamconv.ItemReader) (dst streamconv.ItemReader, err error) {
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

	var extractor windowExtractor
	if t.skipPartial {
		extractor = &completeWindowExtractor{
			reader:  streamconv.NewLookaheadItemReader(src),
			size:    t.size,
			step:    t.step,
			overlap: overlap,
			hole:    hole,
		}
	} else {
		extractor = &partialWindowExtractor{
			reader:     streamconv.NewLookaheadItemReader(src),
			size:       t.size,
			step:       t.step,
			overlap:    overlap,
			hole:       hole,
			bufferStep: bufferStep,
		}
	}

	dst = &multiWindowReader{
		extractor:  extractor,
		subProgram: t.subProgram,
	}
	return
}

func NewWindowTransformer(subProgram streamconv.Transformer, size uint, step uint, skipPartial bool) streamconv.Transformer {
	return &windowTransformer{
		subProgram:  subProgram,
		size:        size,
		step:        step,
		skipPartial: skipPartial,
	}
}
