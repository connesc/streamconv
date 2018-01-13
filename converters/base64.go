package converters

import (
	"bytes"
	"encoding/base64"
	"streamconv"
)

type encoder struct {
	buffer *bytes.Buffer
}

func (c *encoder) Convert(src []byte) (dst []byte, err error) {
	c.buffer.Reset()
	encoder := base64.NewEncoder(base64.StdEncoding, c.buffer)
	_, err = encoder.Write(src)
	if err == nil {
		err = encoder.Close()
	}
	return c.buffer.Bytes(), err
}

func NewBase64Encode() streamconv.Converter {
	return &encoder{&bytes.Buffer{}}
}

type decoder struct {
	buffer *bytes.Buffer
}

func (c *decoder) Convert(src []byte) (dst []byte, err error) {
	c.buffer.Reset()
	decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(src))
	_, err = c.buffer.ReadFrom(decoder)
	return c.buffer.Bytes(), err
}

func NewBase64Decode() streamconv.Converter {
	return &decoder{&bytes.Buffer{}}
}
