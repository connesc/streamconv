package extractors

import (
	"bufio"
	"bytes"
	"io"
	"regexp"

	"github.com/connesc/streamconv"
)

type splitExtractor struct {
	scanner *bufio.Scanner
}

// TODO: handle errors

func (r *splitExtractor) ReadItem() (item io.Reader, err error) {
	if !r.scanner.Scan() {
		return nil, io.EOF
	}

	return bytes.NewReader(r.scanner.Bytes()), r.scanner.Err()
}

func NewSplitExtractor(in io.Reader, delimiter *regexp.Regexp) streamconv.ItemReader {
	scanner := bufio.NewScanner(in)

	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		indices := delimiter.FindIndex(data)
		if indices != nil {
			advance = indices[1]
			token = data[0:indices[0]]
		} else if atEOF && len(data) > 0 {
			advance = len(data)
			token = data
		}
		return
	})

	return &splitExtractor{
		scanner: scanner,
	}
}
