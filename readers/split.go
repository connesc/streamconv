package readers

import (
	"bufio"
	"bytes"
	"io"
	"streamconv"
)

type splitReader struct {
	scanner *bufio.Scanner
}

// TODO: improve streaming and handle custom delimiter

func (r *splitReader) ReadItem() (item io.Reader, err error) {
	if !r.scanner.Scan() {
		return nil, io.EOF
	}

	return bytes.NewReader(r.scanner.Bytes()), r.scanner.Err()
}

func NewSplitReader(in io.Reader) streamconv.ItemReader {
	return &splitReader{
		scanner: bufio.NewScanner(in),
	}
}
