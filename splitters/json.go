package splitters

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/connesc/streamconv"
)

type jsonSplitter struct {
	decoder *json.Decoder
	buffer  json.RawMessage
}

func (r *jsonSplitter) ReadItem() (item io.Reader, err error) {
	err = r.decoder.Decode(&r.buffer)
	return bytes.NewReader(r.buffer), err
}

func NewJSONSplitter(in io.Reader) streamconv.Splitter {
	return &jsonSplitter{
		decoder: json.NewDecoder(in),
	}
}
