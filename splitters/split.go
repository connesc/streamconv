package splitters

import (
	"bufio"
	"bytes"
	"io"
	"regexp"
	"streamconv"
)

type simpleSplitter struct {
	scanner *bufio.Scanner
}

// TODO: handle errors

func (r *simpleSplitter) ReadItem() (item io.Reader, err error) {
	if !r.scanner.Scan() {
		return nil, io.EOF
	}

	return bytes.NewReader(r.scanner.Bytes()), r.scanner.Err()
}

func NewSimpleSplitter(in io.Reader, delim string) streamconv.Splitter {
	scanner := bufio.NewScanner(in)
	re, _ := regexp.Compile(delim) // TODO: handle error

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		indices := re.FindIndex(data)
		if indices != nil {
			advance = indices[1]
			token = data[0:indices[0]]
		} else if atEOF && len(data) > 0 {
			advance = len(data)
			token = data
		}
		return
	})

	return &simpleSplitter{
		scanner: scanner,
	}
}
