package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type varintExtractorCommand struct {
	name string
}

func (c *varintExtractorCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *varintExtractorCommand) Parse(args []string, in io.Reader) (extractor streamconv.ItemReader, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewVarintExtractor(in), nil
}

func NewVarintExtractorCommand(name string) (command streamconv.ExtractorCommand) {
	return &varintExtractorCommand{name}
}

func RegisterVarintExtractor(name string) {
	streamconv.RegisterExtractor(name, NewVarintExtractorCommand(name))
}
