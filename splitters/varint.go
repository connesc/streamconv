package splitters

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"streamconv"

	"github.com/golang/protobuf/proto"
)

const maxVarintSize = 10 // math.Ceil(64 / 7)

type fixedReader struct {
	in   io.Reader
	size int
}

func (r *fixedReader) Read(p []byte) (n int, err error) {
	if r.size <= 0 {
		return 0, io.EOF
	}
	if len(p) > r.size {
		p = p[0:r.size]
	}
	n, err = r.in.Read(p)
	r.size -= n
	if err == io.EOF && r.size != 0 {
		err = io.ErrUnexpectedEOF
	}
	return
}

type varintSplitter struct {
	in     *bufio.Reader
	varint *proto.Buffer
}

func (r *varintSplitter) ReadItem() (item io.Reader, err error) {
	head, err := r.in.Peek(maxVarintSize)
	if (err == io.EOF && len(head) == 0) || (err != nil && err != io.EOF) {
		return
	}

	r.varint.SetBuf(head)
	size, err := r.varint.DecodeVarint()
	if err != nil {
		return
	}

	if size > math.MaxInt32 {
		return nil, fmt.Errorf("invalid item size: %v", size)
	}

	_, err = r.in.Discard(proto.SizeVarint(size))
	if err != nil {
		return
	}

	return &fixedReader{in: r.in, size: int(size)}, nil
}

func NewVarintSplitter(in io.Reader) streamconv.Splitter {
	return &varintSplitter{
		in:     bufio.NewReaderSize(in, maxVarintSize),
		varint: &proto.Buffer{},
	}
}
