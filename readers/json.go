package readers

import (
	"bytes"
	"encoding/json"
	"io"
	"streamconv"
)

type jsonReader struct {
	decoder *json.Decoder
	buffer  json.RawMessage
}

func (r *jsonReader) ReadItem() (item io.Reader, err error) {
	err = r.decoder.Decode(&r.buffer)
	return bytes.NewReader(r.buffer), err
}

func NewJSONReader(in io.Reader) streamconv.ItemReader {
	return &jsonReader{
		decoder: json.NewDecoder(in),
	}
}
