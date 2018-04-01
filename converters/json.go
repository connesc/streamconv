package converters

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/connesc/streamconv"
)

type indenter struct {
	prefix    string
	indent    string
	srcBuffer *bytes.Buffer
	dstBuffer *bytes.Buffer
}

func (c *indenter) Convert(src io.Reader) (dst io.Reader, err error) {
	c.srcBuffer.Reset()
	_, err = c.srcBuffer.ReadFrom(src)
	if err != nil {
		return
	}

	c.dstBuffer.Reset()
	err = json.Indent(c.dstBuffer, c.srcBuffer.Bytes(), c.prefix, c.indent)
	if err != nil {
		return
	}

	return c.dstBuffer, nil
}

func NewJSONIndenter(prefix string, indent string) streamconv.Converter {
	return &indenter{prefix, indent, &bytes.Buffer{}, &bytes.Buffer{}}
}

type compactor struct {
	srcBuffer *bytes.Buffer
	dstBuffer *bytes.Buffer
}

func (c *compactor) Convert(src io.Reader) (dst io.Reader, err error) {
	c.srcBuffer.Reset()
	_, err = c.srcBuffer.ReadFrom(src)
	if err != nil {
		return
	}

	c.dstBuffer.Reset()
	err = json.Compact(c.dstBuffer, c.srcBuffer.Bytes())
	if err != nil {
		return
	}

	return c.dstBuffer, nil
}

func NewJSONCompactor() streamconv.Converter {
	return &compactor{&bytes.Buffer{}, &bytes.Buffer{}}
}
