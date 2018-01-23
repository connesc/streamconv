package extractors

import (
	"io"

	"github.com/connesc/streamconv"
)

type singleExtractor struct {
	in   io.Reader
	done bool
}

func (r *singleExtractor) ReadItem() (item io.Reader, err error) {
	if r.done {
		return nil, io.EOF
	}

	r.done = true
	return r.in, nil
}

func NewSingleExtractor(in io.Reader) streamconv.ItemReader {
	return &singleExtractor{
		in:   in,
		done: false,
	}
}
