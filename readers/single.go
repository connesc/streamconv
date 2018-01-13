package readers

import (
	"bytes"
	"io"
	"streamconv"
)

type singleReader struct {
	in   io.Reader
	done bool
}

func (r *singleReader) ReadItem() (item []byte, err error) {
	if r.done {
		return nil, io.EOF
	}

	buffer := bytes.Buffer{}
	_, err = buffer.ReadFrom(r.in)
	r.done = true
	return buffer.Bytes(), err
}

func NewSingleReader(in io.Reader) streamconv.ItemReader {
	return &singleReader{in, false}
}
