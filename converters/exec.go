package converters

import (
	"io"
	"os"
	"os/exec"

	"github.com/connesc/streamconv"
)

type executor struct {
	name string
	args []string
}

type blockingReader struct {
	reader io.Reader
	wait   func() error
}

func (r *blockingReader) Read(p []byte) (n int, err error) {
	n, err = r.reader.Read(p)
	if err != nil {
		if initialErr := r.wait(); initialErr != nil {
			err = initialErr
		}
	}
	return
}

func (c *executor) Convert(src io.Reader) (dst io.Reader, err error) {
	cmd := exec.Command(c.name, c.args...)
	cmd.Stdin = src
	cmd.Stderr = os.Stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	dst = &blockingReader{stdout, cmd.Wait}
	err = cmd.Start()
	return
}

func NewExecutor(name string, args []string) streamconv.Converter {
	return &executor{
		name: name,
		args: args,
	}
}
