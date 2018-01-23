package streamconv

import (
	"fmt"
	"io"
)

var commands = map[string]Command{}
var splitters = map[string]SplitterCommand{}
var converters = map[string]ConverterCommand{}
var joiners = map[string]JoinerCommand{}

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

func RegisterSplitter(name string, splitter SplitterCommand) {
	commands[name] = splitter
	splitters[name] = splitter
}

func RegisterConverter(name string, converter ConverterCommand) {
	commands[name] = converter
	converters[name] = converter
}

func RegisterJoiner(name string, joiner JoinerCommand) {
	commands[name] = joiner
	joiners[name] = joiner
}

func GetSplitter(tokens []string, in io.Reader) (splitter Splitter, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty splitter command")
	}
	command, ok := splitters[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("splitter \"%v\" cannot be found", tokens[0])
	}
	return command.Parse(tokens[1:], in)
}

func ApplyConverter(tokens []string, input Splitter) (output Splitter, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty converter command")
	}

	name := tokens[0]
	args := tokens[1:]

	if len(name) > 0 && name[0] == '*' {
		name = name[1:]
		command, ok := splitters[name]
		if !ok {
			return nil, fmt.Errorf("splitter \"%v\" cannot be found", name)
		}
		output = &spreadReader{
			reader:          input,
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
		output = &converterReader{
			reader:    input,
			converter: converter,
		}
	}

	return
}

func GetJoiner(tokens []string, out io.Writer) (joiner Joiner, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty joiner command")
	}
	command, ok := joiners[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("joiner \"%v\" cannot be found", tokens[0])
	}
	return command.Parse(tokens[1:], out)
}

func PrintUsage(name string, output io.Writer) (err error) {
	command, ok := commands[name]
	if !ok {
		return fmt.Errorf("command \"%v\" cannot be found", name)
	}
	err = command.PrintUsage(output)
	return
}

type converterReader struct {
	reader    Splitter
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
	reader          Splitter
	spreaderCommand SplitterCommand
	spreaderArgs    []string
	spreader        Splitter
}

func (r *spreadReader) ReadItem() (item io.Reader, err error) {
	if r.spreader == nil {
		err = io.EOF
	} else {
		item, err = r.spreader.ReadItem()
	}

	for err == io.EOF {
		var input io.Reader
		input, err = r.reader.ReadItem()
		if err != nil {
			return
		}

		r.spreader, err = r.spreaderCommand.Parse(r.spreaderArgs, input)
		if err != nil {
			return
		}

		item, err = r.spreader.ReadItem()
	}

	return
}
