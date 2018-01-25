package streamconv

import (
	"fmt"
	"io"
)

var clis = map[string]CLI{}
var extractors = map[string]ExtractorCLI{}
var converters = map[string]ConverterCLI{}
var combiners = map[string]CombinerCLI{}

func getNames(clis map[string]CLI) (names []string) {
	names = make([]string, len(clis))
	index := 0
	for name := range clis {
		names[index] = name
		index++
	}
	return
}

func Commands() (names []string) {
	return getNames(clis)
}

func RegisterExtractor(name string, extractor ExtractorCLI) {
	clis[name] = extractor
	extractors[name] = extractor
}

func RegisterConverter(name string, converter ConverterCLI) {
	clis[name] = converter
	converters[name] = converter
}

func RegisterCombiner(name string, combiner CombinerCLI) {
	clis[name] = combiner
	combiners[name] = combiner
}

func GetExtractor(tokens []string, in io.Reader) (extractor ItemReader, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty extractor command")
	}
	cli, ok := extractors[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("extractor \"%v\" cannot be found", tokens[0])
	}
	command, err := cli.Parse(tokens[1:])
	if err != nil {
		return
	}
	return command.Run(in)
}

func ApplyConverter(tokens []string, in ItemReader) (out ItemReader, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty converter command")
	}

	name := tokens[0]
	args := tokens[1:]

	if len(name) > 0 && name[0] == '*' {
		name = name[1:]
		cli, ok := extractors[name]
		if !ok {
			return nil, fmt.Errorf("extractor \"%v\" cannot be found", name)
		}
		command, err := cli.Parse(args)
		if err != nil {
			return nil, err
		}
		out = &spreadReader{
			reader:  in,
			command: command,
		}
	} else {
		cli, ok := converters[name]
		if !ok {
			return nil, fmt.Errorf("converter \"%v\" cannot be found", name)
		}
		command, err := cli.Parse(args)
		if err != nil {
			return nil, err
		}
		converter, err := command.Run()
		if err != nil {
			return nil, err
		}
		out = &converterReader{
			reader:    in,
			converter: converter,
		}
	}

	return
}

func GetCombiner(tokens []string, out io.Writer) (combiner ItemWriter, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty combiner command")
	}
	cli, ok := combiners[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("combiner \"%v\" cannot be found", tokens[0])
	}
	command, err := cli.Parse(tokens[1:])
	if err != nil {
		return
	}
	return command.Run(out)
}

func PrintUsage(name string, out io.Writer) (err error) {
	cli, ok := clis[name]
	if !ok {
		return fmt.Errorf("command \"%v\" cannot be found", name)
	}
	err = cli.PrintUsage(out)
	return
}

type converterReader struct {
	reader    ItemReader
	converter Converter
}

func (r *converterReader) ReadItem() (item io.Reader, err error) {
	item, err = r.reader.ReadItem()
	if err != nil {
		return nil, err
	}

	item, err = r.converter.Convert(item)
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	return
}

type spreadReader struct {
	reader   ItemReader
	command  ExtractorCommand
	spreader ItemReader
}

func (r *spreadReader) ReadItem() (item io.Reader, err error) {
	if r.spreader == nil {
		err = io.EOF
	} else {
		item, err = r.spreader.ReadItem()
	}

	for err == io.EOF {
		var src io.Reader
		src, err = r.reader.ReadItem()
		if err != nil {
			return
		}

		r.spreader, err = r.command.Run(src)
		if err != nil {
			return
		}

		item, err = r.spreader.ReadItem()
	}

	return
}
