package splitters

import (
	"bufio"
	"bytes"
	"io"
	"streamconv"
)

type windowSplitter struct {
	in      *bufio.Reader
	size    int
	step    int
	partial bool
	started bool
	done    bool
}

func (r *windowSplitter) ReadItem() (item io.Reader, err error) {
	if r.done {
		return nil, io.EOF
	}

	if r.started {
		_, err = r.in.Discard(r.step)
		if err != nil {
			if err == io.EOF {
				r.done = true
			}
			return
		}
	}

	window, err := r.in.Peek(r.size)
	if err != nil {
		if err != io.EOF {
			return
		}
		r.done = true
		if !r.partial || len(window) == 0 || (r.step < r.size && r.started && len(window) == (r.size-r.step)) {
			return nil, io.EOF
		}
	}

	r.started = true
	return bytes.NewReader(window), nil
}

func NewWindowSplitter(in io.Reader, size int, step int, partial bool) streamconv.Splitter {
	return &windowSplitter{
		in:      bufio.NewReaderSize(in, size),
		size:    size,
		step:    step,
		partial: partial,
		started: false,
		done:    false,
	}
}
