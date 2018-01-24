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
	done   <-chan error
}

func (r *blockingReader) checkError() (err error) {
	select {
	case err = <-r.done:
	default:
	}
	return
}

func (r *blockingReader) Read(p []byte) (n int, err error) {
	err = r.checkError()
	if err != nil {
		return
	}

	n, err = r.reader.Read(p)

	if err == io.EOF {
		err = <-r.done
		if err == nil {
			err = io.EOF
		}
	} else if err == nil {
		err = r.checkError()
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

	done := make(chan error)
	dst = &blockingReader{stdout, done}

	err = cmd.Start()
	if err != nil {
		return
	}

	go func() {
		done <- cmd.Wait()
	}()

	return
}

func NewExecutor(name string, args []string) streamconv.Converter {
	return &executor{
		name: name,
		args: args,
	}
}
