package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type varintExtractorCLI struct {
	name string
}

func (c *varintExtractorCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *varintExtractorCLI) Parse(args []string) (command streamconv.ExtractorCommand, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func(in io.Reader) (streamconv.ItemReader, error) {
		return NewVarintExtractor(in), nil
	}
	return
}

func NewVarintExtractorCLI(name string) (cli streamconv.ExtractorCLI) {
	return &varintExtractorCLI{name}
}

func RegisterVarintExtractor(name string) {
	streamconv.RegisterExtractor(name, NewVarintExtractorCLI(name))
}
