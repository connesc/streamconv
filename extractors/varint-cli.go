package extractors

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"
)

type varintExtractorCommand struct{}

func (c *varintExtractorCommand) Run(in io.Reader) (extractor streamconv.ItemReader, err error) {
	return NewVarintExtractor(in), nil
}

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

	return &varintExtractorCommand{}, nil
}

func NewVarintExtractorCLI(name string) (cli streamconv.ExtractorCLI) {
	return &varintExtractorCLI{name}
}

func RegisterVarintExtractor(name string) {
	streamconv.RegisterExtractor(name, NewVarintExtractorCLI(name))
}
