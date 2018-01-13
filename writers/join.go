package writers

import (
	"io"
	"streamconv"
	"strings"
)

type joinWriter struct {
	out     io.Writer
	delim   string
	started bool
}

func (w *joinWriter) WriteItem(item []byte) (err error) {
	if w.started && len(w.delim) > 0 {
		_, err = strings.NewReader(w.delim).WriteTo(w.out)
		if err != nil {
			return
		}
	}

	_, err = w.out.Write(item)
	w.started = true
	return
}

func NewJoinWriter(out io.Writer, delim string) streamconv.ItemWriter {
	return &joinWriter{out, delim, false}
}
