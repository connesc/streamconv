package combiners

import (
	"bytes"
	"io"

	"github.com/connesc/streamconv"

	"github.com/golang/protobuf/proto"
)

type varintCombiner struct {
	out    io.Writer
	varint *proto.Buffer
	buffer *bytes.Buffer
}

func (w *varintCombiner) WriteItem(item io.Reader) (err error) {
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

func NewVarintCombiner(out io.Writer) streamconv.ItemWriter {
	return &varintCombiner{
		out:    out,
		varint: &proto.Buffer{},
		buffer: &bytes.Buffer{},
	}
}
