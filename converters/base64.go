package converters

import (
	"encoding/base64"
	"streamconv"
)

// TODO: reuse buffers

type encoder struct{}

func (encoder) Convert(src []byte) (dst []byte, err error) {
	dst = make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst, nil
}

func NewBase64Encode() streamconv.Converter {
	return &encoder{}
}

type decoder struct{}

func (decoder) Convert(src []byte) (dst []byte, err error) {
	dst = make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	n, err := base64.StdEncoding.Decode(dst, src)
	return dst[:n], err
}

func NewBase64Decode() streamconv.Converter {
	return &decoder{}
}
