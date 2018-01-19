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

func GetConverter(tokens []string) (converter Converter, err error) {
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty converter command")
	}
	command, ok := converters[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("converter \"%v\" cannot be found", tokens[0])
	}
	return command.Parse(tokens[1:])
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
