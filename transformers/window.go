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

type windowReader interface {
	streamconv.ItemReader
	Initialize(first bool) (err error)
	Finalize() (err error)
}

type partialWindowReader struct {
	reader streamconv.LookaheadItemReader

	size       uint
	step       uint
	overlap    uint
	hole       uint
	bufferStep uint

	buffer        *windowBuffer
	bufferIndex   uint
	reusedCount   uint
	readCount     uint
	bufferedCount uint
	err           error
}

func newPartialWindowReader(reader streamconv.ItemReader, size uint, step uint, overlap uint, hole uint) windowReader {
	var bufferStep uint
	if overlap > step {
		bufferStep = step
	} else {
		bufferStep = overlap
	}

	return &partialWindowReader{
		reader: streamconv.NewLookaheadItemReader(reader),

		size:       size,
		step:       step,
		overlap:    overlap,
		hole:       hole,
		bufferStep: bufferStep,

		buffer:        newWindowBuffer(overlap),
		bufferIndex:   0,
		reusedCount:   0,
		readCount:     size,
		bufferedCount: overlap,
	}
}

func (r *partialWindowReader) Initialize(first bool) (err error) {
	if !first {
		err = skip(r.reader, r.hole)
		if err != nil {
			return
		}
		r.bufferIndex -= r.overlap
		r.reusedCount = r.overlap
		r.readCount = r.size - r.overlap
		r.bufferedCount = r.bufferStep
		r.err = nil
	}
	_, err = r.reader.Lookahead()
	return
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

func (r *partialWindowReader) Finalize() (err error) {
	return r.err
}

type completeWindowReader struct {
	reader streamconv.ItemReader

	size        uint
	step        uint
	overlap     uint
	hole        uint
	forwardLast bool

	buffer      *windowBuffer
	bufferIndex uint
	lastItem    io.Reader
	index       uint
	err         error
}

func newCompleteWindowReader(reader streamconv.ItemReader, size uint, step uint, overlap uint, hole uint) windowReader {
	forwardLast := overlap == 0
	bufferLength := size
	if forwardLast {
		bufferLength--
	}

	return &completeWindowReader{
		reader: reader,

		size:        size,
		step:        step,
		overlap:     overlap,
		hole:        hole,
		forwardLast: forwardLast,

		buffer:      newWindowBuffer(bufferLength),
		bufferIndex: 0,
		index:       0,
	}
}

func (r *completeWindowReader) Initialize(first bool) (err error) {
	if !first {
		err = skip(r.reader, r.hole)
		if err != nil {
			return
		}
		r.bufferIndex -= r.overlap
		r.index = r.overlap
		r.err = nil
	}

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

func (r *completeWindowReader) Finalize() (err error) {
	return r.err
}

type multiWindowReader struct {
	reader      windowReader
	subProgram  streamconv.Transformer
	first       bool
	transformed streamconv.ItemReader
	err         error
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
		if r.transformed == nil {
			if !r.first {
				err = r.reader.Finalize()
				if err != io.EOF {
					if err == nil {
						err = fmt.Errorf("the previous window has not been fully read")
					}
					return
				}
			}

			err = r.reader.Initialize(r.first)
			if err != nil {
				return
			}

			r.first = false

			if r.subProgram != nil {
				r.transformed, err = r.subProgram.Transform(r.reader)
				if err != nil {
					if err == io.EOF {
						err = io.ErrUnexpectedEOF
					}
					return
				}
			} else {
				r.transformed = r.reader
			}
		}

		item, err = r.transformed.ReadItem()
		if err != io.EOF {
			return
		}
		err = nil
		r.transformed = nil
	}
}

type windowTransformer struct {
	subProgram  streamconv.Transformer
	size        uint
	step        uint
	skipPartial bool
}

func (t *windowTransformer) Transform(src streamconv.ItemReader) (dst streamconv.ItemReader, err error) {
	var overlap, hole uint
	if t.step < t.size {
		overlap = t.size - t.step
	} else {
		hole = t.step - t.size
	}

	var reader windowReader
	if t.skipPartial {
		reader = newCompleteWindowReader(src, t.size, t.step, overlap, hole)
	} else {
		reader = newPartialWindowReader(src, t.size, t.step, overlap, hole)
	}

	dst = &multiWindowReader{
		reader:     reader,
		subProgram: t.subProgram,
		first:      true,
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
