package converters

import (
	"bytes"
	"encoding/base64"
	"streamconv"
)

type encoder struct {
	encoding *base64.Encoding
	buffer   *bytes.Buffer
}

func (c *encoder) Convert(src []byte) (dst []byte, err error) {
	c.buffer.Reset()
	c.buffer.Grow(c.encoding.EncodedLen(len(src)))
	encoder := base64.NewEncoder(c.encoding, c.buffer)
	_, err = encoder.Write(src)
	if err == nil {
		err = encoder.Close()
	}
	return c.buffer.Bytes(), err
}

func NewBase64Encode() streamconv.Converter {
	return &encoder{base64.StdEncoding, &bytes.Buffer{}}
}

type decoder struct {
	encoding *base64.Encoding
	buffer   *bytes.Buffer
}

func (c *decoder) Convert(src []byte) (dst []byte, err error) {
	c.buffer.Reset()
	c.buffer.Grow(c.encoding.DecodedLen(len(src)))
	decoder := base64.NewDecoder(c.encoding, bytes.NewReader(src))
	_, err = c.buffer.ReadFrom(decoder)
	return c.buffer.Bytes(), err
}

func NewBase64Decode() streamconv.Converter {
	return &decoder{base64.StdEncoding, &bytes.Buffer{}}
}
