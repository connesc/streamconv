package writers

import (
	"bytes"
	"io"
	"streamconv"

	"github.com/golang/protobuf/proto"
)

type varintWriter struct {
	out    io.Writer
	buffer *proto.Buffer
}

func (w *varintWriter) WriteItem(item []byte) (err error) {
	w.buffer.Reset()
	err = w.buffer.EncodeVarint(uint64(len(item)))
	if err != nil {
		return
	}

	_, err = w.out.Write(w.buffer.Bytes())
	if err != nil {
		return
	}

	_, err = bytes.NewReader(item).WriteTo(w.out)
	return
}

func NewVarintWriter(out io.Writer) streamconv.ItemWriter {
	return &varintWriter{out, &proto.Buffer{}}
}
