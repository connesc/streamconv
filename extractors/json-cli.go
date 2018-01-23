package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type jsonExtractorCommand struct {
	name string
}

func (c *jsonExtractorCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *jsonExtractorCommand) Parse(args []string, in io.Reader) (extractor streamconv.ItemReader, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewJSONExtractor(in), nil
}

func NewJSONExtractorCommand(name string) (command streamconv.ExtractorCommand) {
	return &jsonExtractorCommand{name}
}

func RegisterJSONExtractor(name string) {
	streamconv.RegisterExtractor(name, NewJSONExtractorCommand(name))
}
