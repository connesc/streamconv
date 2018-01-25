package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type jsonExtractorCommand struct{}

func (c *jsonExtractorCommand) Run(in io.Reader) (extractor streamconv.ItemReader, err error) {
	return NewJSONExtractor(in), nil
}

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

	return &jsonExtractorCommand{}, nil
}

func NewJSONExtractorCLI(name string) (cli streamconv.ExtractorCLI) {
	return &jsonExtractorCLI{name}
}

func RegisterJSONExtractor(name string) {
	streamconv.RegisterExtractor(name, NewJSONExtractorCLI(name))
}
