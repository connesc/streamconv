package app

import (
	"fmt"
	"io"
	"os"

	"github.com/connesc/streamconv"
)

type App struct {
	splitter   streamconv.Splitter
	converters []streamconv.Converter
	joiner     streamconv.Joiner
}

func New(program string) (app *App, err error) {
	commands, err := parse(program)
	if err != nil {
		return
	}

	if len(commands) < 2 {
		return nil, fmt.Errorf("not enough commands")
	}

	splitter, err := streamconv.GetSplitter(commands[0], os.Stdin)
	if err != nil {
		return
	}

	converters := make([]streamconv.Converter, len(commands)-2)
	for index, command := range commands[1 : len(commands)-1] {
		converters[index], err = streamconv.GetConverter(command)
		if err != nil {
			return
		}
	}

	joiner, err := streamconv.GetJoiner(commands[len(commands)-1], os.Stdout)
	if err != nil {
		return
	}

	app = &App{splitter, converters, joiner}
	return
}

func (app *App) Run() (err error) {
	for {
		item, err := app.splitter.ReadItem()
		if err != nil {
			return err
		}

		for _, converter := range app.converters {
			item, err = converter.Convert(item)
			if err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				return err
			}
		}

		err = app.joiner.WriteItem(item)
		if err != nil {
			return err
		}
	}
}
