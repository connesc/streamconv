package converters

import (
	"io"
	"os/exec"

	"github.com/connesc/streamconv"
)

type executor struct {
	name string
	args []string
}

func (c *executor) Convert(src io.Reader) (dst io.Reader, err error) {
	cmd := exec.Command(c.name, c.args...)
	cmd.Stdin = src
	dst, err = cmd.StdoutPipe()
	if err != nil {
		return
	}

	err = cmd.Start()
	return

	// TODO: handle exit status and maybe Stderr
}

func NewExecutor(name string, args []string) streamconv.Converter {
	return &executor{
		name: name,
		args: args,
	}
}
