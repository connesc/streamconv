package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type singleExtractorCommand struct {
	name string
}

func (c *singleExtractorCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *singleExtractorCommand) Parse(args []string, in io.Reader) (extractor streamconv.ItemReader, err error) {
	if len(args) > 0 {
		return nil, fmt.Errorf("too many arguments (expected 0, got %v)", len(args))
	}

	return NewSingleExtractor(in), nil
}

func NewSingleExtractorCommand(name string) (command streamconv.ExtractorCommand) {
	return &singleExtractorCommand{name}
}

func RegisterSingleExtractor(name string) {
	streamconv.RegisterExtractor(name, NewSingleExtractorCommand(name))
}
