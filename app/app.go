package app

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type App interface {
	Run(dst io.Writer, src io.Reader) error
}

type streamconvApp [][]string

func New(program string) (app App, err error) {
	commands, err := parse(program)
	return streamconvApp(commands), err
}

func (app streamconvApp) Run(dst io.Writer, src io.Reader) (err error) {
	if len(app) < 2 {
		return fmt.Errorf("not enough commands")
	}

	reader, err := streamconv.GetExtractor(app[0], src)
	if err != nil {
		return
	}

	for _, command := range app[1 : len(app)-1] {
		reader, err = streamconv.ApplyConverter(command, reader)
		if err != nil {
			return
		}
	}

	writer, err := streamconv.GetCombiner(app[len(app)-1], dst)
	if err != nil {
		return
	}

	for {
		item, err := reader.ReadItem()
		if err != nil {
			return err
		}

		err = writer.WriteItem(item)
		if err != nil {
			return err
		}
	}
}
