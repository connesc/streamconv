package readers

import (
	"io"
	"streamconv"
)

type singleReader struct {
	in   io.Reader
	done bool
}

func (r *singleReader) ReadItem() (item io.Reader, err error) {
	if r.done {
		return nil, io.EOF
	}

	r.done = true
	return r.in, nil
}

func NewSingleReader(in io.Reader) streamconv.ItemReader {
	return &singleReader{
		in:   in,
		done: false,
	}
}
