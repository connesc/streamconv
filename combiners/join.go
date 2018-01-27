package combiners

import (
	"io"
	"strings"

	"github.com/connesc/streamconv"
)

type joinCombiner struct {
	out       io.Writer
	delimiter string
	started   bool
}

func (w *joinCombiner) WriteItem(item io.Reader) (err error) {
	if w.started && len(w.delimiter) > 0 {
		_, err = strings.NewReader(w.delimiter).WriteTo(w.out)
		if err != nil {
			return
		}
	}

	_, err = io.Copy(w.out, item)
	w.started = true
	return
}

func NewJoinCombiner(out io.Writer, delimiter string) streamconv.ItemWriter {
	return &joinCombiner{
		out:       out,
		delimiter: delimiter,
		started:   false,
	}
}
