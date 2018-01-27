package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type jsonExtractorCLI struct {
	name string
}

func (c *jsonExtractorCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *jsonExtractorCLI) Parse(args []string) (command streamconv.ExtractorCommand, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	command = func(in io.Reader) (streamconv.ItemReader, error) {
		return NewJSONExtractor(in), nil
	}
	return
}

func NewJSONExtractorCLI(name string) (cli streamconv.ExtractorCLI) {
	return &jsonExtractorCLI{name}
}

func RegisterJSONExtractor(name string) {
	streamconv.RegisterExtractor(name, NewJSONExtractorCLI(name))
}
