package converters

import (
	"encoding/base64"
	"io"

	"github.com/connesc/streamconv"
)

type encoder struct {
	encoding *base64.Encoding
}

func (c *encoder) Convert(src io.Reader) (dst io.Reader, err error) {
	pr, pw := io.Pipe()
	encoder := base64.NewEncoder(c.encoding, pw)

	go func() {
		_, err := io.Copy(encoder, src)
		if err == nil {
			err = encoder.Close()
		}
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewBase64Encoder() streamconv.Converter {
	return &encoder{base64.StdEncoding}
}

type decoder struct {
	encoding *base64.Encoding
}

func (c *decoder) Convert(src io.Reader) (dst io.Reader, err error) {
	return base64.NewDecoder(c.encoding, src), nil
}

func NewBase64Decoder() streamconv.Converter {
	return &decoder{
		encoding: base64.StdEncoding,
	}
}
