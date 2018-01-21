package splitters

import (
	"io"

	"github.com/connesc/streamconv"
)

type singleSplitter struct {
	in   io.Reader
	done bool
}

func (r *singleSplitter) ReadItem() (item io.Reader, err error) {
	if r.done {
		return nil, io.EOF
	}

	r.done = true
	return r.in, nil
}

func NewSingleSplitter(in io.Reader) streamconv.Splitter {
	return &singleSplitter{
		in:   in,
		done: false,
	}
}
