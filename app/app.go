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

	splitter, err := streamconv.GetSplitter(app[0], src)
	if err != nil {
		return
	}

	converters := make([]streamconv.Converter, len(app)-2)
	for index, command := range app[1 : len(app)-1] {
		converters[index], err = streamconv.GetConverter(command)
		if err != nil {
			return
		}
	}

	joiner, err := streamconv.GetJoiner(app[len(app)-1], dst)
	if err != nil {
		return
	}

	for {
		item, err := splitter.ReadItem()
		if err != nil {
			return err
		}

		for _, converter := range converters {
			item, err = converter.Convert(item)
			if err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				return err
			}
		}

		err = joiner.WriteItem(item)
		if err != nil {
			return err
		}
	}
}
