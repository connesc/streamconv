package readers

import (
	"encoding/json"
	"io"
	"streamconv"
)

type jsonReader struct {
	decoder *json.Decoder
}

func (r *jsonReader) ReadItem() (item []byte, err error) {
	err = r.decoder.Decode((*json.RawMessage)(&item))
	return
}

func NewJSONReader(in io.Reader) streamconv.ItemReader {
	return &jsonReader{json.NewDecoder(in)}
}
