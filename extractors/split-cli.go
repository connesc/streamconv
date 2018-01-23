package extractors

import (
	"fmt"
	"io"
	"regexp"

	"github.com/connesc/streamconv"
)

type splitExtractorCommand struct {
	name string
}

func (c *splitExtractorCommand) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *splitExtractorCommand) Parse(args []string, in io.Reader) (extractor streamconv.ItemReader, err error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("too many arguments (expected up to 1, got %v)", len(args))
	}

	rawDelimiter := "\\n"
	if len(args) == 1 {
		rawDelimiter = args[0]
	}

	delimiter, err := regexp.Compile(rawDelimiter)
	if err != nil {
		return
	}

	return NewSplitExtractor(in, delimiter), nil
}

func NewSplitExtractorCommand(name string) (command streamconv.ExtractorCommand) {
	return &splitExtractorCommand{name}
}

func RegisterSplitExtractor(name string) {
	streamconv.RegisterExtractor(name, NewSplitExtractorCommand(name))
}
