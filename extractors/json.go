package extractors

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/connesc/streamconv"
)

type jsonExtractor struct {
	decoder *json.Decoder
	buffer  json.RawMessage
}

func (r *jsonExtractor) ReadItem() (item io.Reader, err error) {
	err = r.decoder.Decode(&r.buffer)
	return bytes.NewReader(r.buffer), err
}

func NewJSONExtractor(in io.Reader) streamconv.ItemReader {
	return &jsonExtractor{
		decoder: json.NewDecoder(in),
	}
}
