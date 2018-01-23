package streamconv

import (
	"fmt"
	"io"
)

var commands = map[string]Command{}
var extractors = map[string]ExtractorCommand{}
var converters = map[string]ConverterCommand{}
var combiners = map[string]CombinerCommand{}

func getNames(commands map[string]Command) (names []string) {
	names = make([]string, len(commands))
	index := 0
	for name := range commands {
		names[index] = name
		index++
	}
	return
}

func Commands() (names []string) {
	return getNames(commands)
}

func RegisterExtractor(name string, extractor ExtractorCommand) {
	commands[name] = extractor
	extractors[name] = extractor
}

func RegisterConverter(name string, converter ConverterCommand) {
	commands[name] = converter
	converters[name] = converter
}

func RegisterCombiner(name string, combiner CombinerCommand) {
	commands[name] = combiner
	combiners[name] = combiner
}

func GetExtractor(tokens []string, in io.Reader) (extractor ItemReader, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty extractor command")
	}
	command, ok := extractors[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("extractor \"%v\" cannot be found", tokens[0])
	}
	return command.Parse(tokens[1:], in)
}

func ApplyConverter(tokens []string, in ItemReader) (out ItemReader, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty converter command")
	}

	name := tokens[0]
	args := tokens[1:]

	if len(name) > 0 && name[0] == '*' {
		name = name[1:]
		command, ok := extractors[name]
		if !ok {
			return nil, fmt.Errorf("extractor \"%v\" cannot be found", name)
		}
		out = &spreadReader{
			reader:          in,
			spreaderCommand: command,
			spreaderArgs:    args,
		}
	} else {
		command, ok := converters[name]
		if !ok {
			return nil, fmt.Errorf("converter \"%v\" cannot be found", name)
		}
		converter, err := command.Parse(args)
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
	command, ok := combiners[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("combiner \"%v\" cannot be found", tokens[0])
	}
	return command.Parse(tokens[1:], out)
}

func PrintUsage(name string, out io.Writer) (err error) {
	command, ok := commands[name]
	if !ok {
		return fmt.Errorf("command \"%v\" cannot be found", name)
	}
	err = command.PrintUsage(out)
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
	reader          ItemReader
	spreaderCommand ExtractorCommand
	spreaderArgs    []string
	spreader        ItemReader
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

		r.spreader, err = r.spreaderCommand.Parse(r.spreaderArgs, src)
		if err != nil {
			return
		}

		item, err = r.spreader.ReadItem()
	}

	return
}
