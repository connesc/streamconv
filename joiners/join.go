package joiners

import (
	"io"
	"strings"

	"github.com/connesc/streamconv"
)

type simpleJoiner struct {
	out     io.Writer
	delim   string
	started bool
}

func (w *simpleJoiner) WriteItem(item io.Reader) (err error) {
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

func NewSimpleJoiner(out io.Writer, delim string) streamconv.Joiner {
	return &simpleJoiner{
		out:     out,
		delim:   delim,
		started: false}
}
