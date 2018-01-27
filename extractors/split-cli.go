package extractors

import (
	"fmt"
	"io"
	"regexp"

	"github.com/connesc/streamconv"
)

type splitExtractorCLI struct {
	name string
}

func (c *splitExtractorCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	return
}

func (c *splitExtractorCLI) Parse(args []string) (command streamconv.ExtractorCommand, err error) {
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

	command = func(in io.Reader) (streamconv.ItemReader, error) {
		return NewSplitExtractor(in, delimiter), nil
	}
	return
}

func NewSplitExtractorCLI(name string) (cli streamconv.ExtractorCLI) {
	return &splitExtractorCLI{name}
}

func RegisterSplitExtractor(name string) {
	streamconv.RegisterExtractor(name, NewSplitExtractorCLI(name))
}
