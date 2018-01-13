package readers

import (
	"bufio"
	"io"
	"streamconv"
)

type splitReader struct {
	scanner *bufio.Scanner
}

// TODO: handle custom delimiter

func (r *splitReader) ReadItem() (item []byte, err error) {
	if !r.scanner.Scan() {
		return nil, io.EOF
	}

	return r.scanner.Bytes(), r.scanner.Err()
}

func NewSplitReader(in io.Reader) streamconv.ItemReader {
	return &splitReader{bufio.NewScanner(in)}
}
