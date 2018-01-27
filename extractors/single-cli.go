package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type singleExtractorCLI struct {
	name string
}

func (c *singleExtractorCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *singleExtractorCLI) Parse(args []string) (command streamconv.ExtractorCommand, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func(in io.Reader) (streamconv.ItemReader, error) {
		return NewSingleExtractor(in), nil
	}
	return
}

func NewSingleExtractorCLI(name string) (cli streamconv.ExtractorCLI) {
	return &singleExtractorCLI{name}
}

func RegisterSingleExtractor(name string) {
	streamconv.RegisterExtractor(name, NewSingleExtractorCLI(name))
}
