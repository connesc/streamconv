package combiners

import (
	"io"
	"strings"

	"github.com/connesc/streamconv"
)

type joinCombiner struct {
	out     io.Writer
	delim   string
	started bool
}

func (w *joinCombiner) WriteItem(item io.Reader) (err error) {
	if w.started && len(w.delim) > 0 {
		_, err = strings.NewReader(w.delim).WriteTo(w.out)
		if err != nil {
			return
		}
	}

	_, err = io.Copy(w.out, item)
	w.started = true
	return
}

func NewJoinCombiner(out io.Writer, delim string) streamconv.ItemWriter {
	return &joinCombiner{
		out:     out,
		delim:   delim,
		started: false}
}
