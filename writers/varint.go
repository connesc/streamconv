package writers

import (
	"bytes"
	"io"
	"streamconv"

	"github.com/golang/protobuf/proto"
)

type varintWriter struct {
	out    io.Writer
	varint *proto.Buffer
	buffer *bytes.Buffer
}

func (w *varintWriter) WriteItem(item io.Reader) (err error) {
	n, err := io.Copy(w.buffer, item)
	if err != nil {
		return
	}

	w.varint.Reset()
	err = w.varint.EncodeVarint(uint64(n))
	if err != nil {
		return
	}

	_, err = w.out.Write(w.varint.Bytes())
	if err != nil {
		return
	}

	_, err = w.buffer.WriteTo(w.out)
	return
}

func NewVarintWriter(out io.Writer) streamconv.ItemWriter {
	return &varintWriter{
		out:    out,
		varint: &proto.Buffer{},
		buffer: &bytes.Buffer{},
	}
}
