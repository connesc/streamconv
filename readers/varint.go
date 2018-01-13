package readers

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"streamconv"

	"github.com/golang/protobuf/proto"
)

const maxVarintSize = 10 // math.Ceil(64 / 7)

type varintReader struct {
	in     *bufio.Reader
	buffer *proto.Buffer
}

func (r *varintReader) ReadItem() (item []byte, err error) {
	head, err := r.in.Peek(maxVarintSize)
	if err != nil {
		return
	}

	r.buffer.SetBuf(head)
	size, err := r.buffer.DecodeVarint()
	if err != nil {
		return
	}

	if size > math.MaxInt32 {
		return nil, fmt.Errorf("invalid item size: %v", size)
	}
	n := int(size)

	_, err = r.in.Discard(proto.SizeVarint(size))
	if err != nil {
		return
	}

	r.in = bufio.NewReaderSize(r.in, n)
	item, err = r.in.Peek(n)
	if err == io.EOF {
		return nil, io.ErrUnexpectedEOF
	} else if err != nil {
		return
	}

	_, err = r.in.Discard(n)
	return
}

func NewVarintReader(in io.Reader) streamconv.ItemReader {
	return &varintReader{bufio.NewReaderSize(in, maxVarintSize), &proto.Buffer{}}
}
