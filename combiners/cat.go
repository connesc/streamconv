package combiners

import (
	"io"

	"github.com/connesc/streamconv"
)

type catCombiner struct {
	out io.Writer
}

func (w *catCombiner) WriteItem(item io.Reader) (err error) {
	_, err = io.Copy(w.out, item)
	return
}

func NewCatCombiner(out io.Writer) streamconv.ItemWriter {
	return &catCombiner{
		out: out,
	}
}
