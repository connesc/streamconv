package joiners

import (
	"bytes"
	"io"
	"streamconv"

	"github.com/golang/protobuf/proto"
)

type varintJoiner struct {
	out    io.Writer
	varint *proto.Buffer
	buffer *bytes.Buffer
}

func (w *varintJoiner) WriteItem(item io.Reader) (err error) {
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

func NewVarintJoiner(out io.Writer) streamconv.Joiner {
	return &varintJoiner{
		out:    out,
		varint: &proto.Buffer{},
		buffer: &bytes.Buffer{},
	}
}
