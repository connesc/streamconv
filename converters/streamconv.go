package converters

import (
	"io"

	"github.com/connesc/streamconv"
	"github.com/connesc/streamconv/app"
)

type streamConverter struct {
	app app.App
}

func (c *streamConverter) Convert(src io.Reader) (dst io.Reader, err error) {
	pr, pw := io.Pipe()

	go func() {
		err := c.app.Run(pw, src)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewStreamConverter(program string) (converter streamconv.Converter, err error) {
	app, err := app.New(program)
	if err != nil {
		return
	}

	converter = &streamConverter{app}
	return
}
