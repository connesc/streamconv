package streamconv

import (
	"fmt"
	"io"
)

var clis = map[string]CLI{}

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

func registerCLI(name string, cli CLI) {
	if _, found := clis[name]; found {
		panic(fmt.Errorf("duplicated command: %v", name))
	}
	clis[name] = cli
}

func RegisterExtractor(name string, extractor ExtractorCLI) {
	registerCLI(name, extractor)
}

func RegisterConverter(name string, converter ConverterCLI) {
	registerCLI(name, converter)
}

func TransformerCombiner(name string, combiner TransformerCLI) {
	registerCLI(name, combiner)
}

func RegisterCombiner(name string, combiner CombinerCLI) {
	registerCLI(name, combiner)
}

func ParseCommand(words []string, subProgram TransformerCommand) (command interface{}, err error) {
	if len(words) == 0 {
		if subProgram != nil {
			return subProgram, nil
		}
		return nil, fmt.Errorf("empty command")
	}

	name := words[0]
	args := words[1:]

	cli, ok := clis[name]
	if !ok {
		return nil, fmt.Errorf("command not found: %v", name)
	}

	switch cli := cli.(type) {
	case ExtractorCLI:
		command, err = cli.Parse(args)
	case ConverterCLI:
		command, err = cli.Parse(args)
	case TransformerCLI:
		command, err = cli.Parse(args, subProgram)
		subProgram = nil
	case CombinerCLI:
		command, err = cli.Parse(args)
	default:
		err = fmt.Errorf("unknown CLI type: %T", cli)
	}

	if err == nil && subProgram != nil {
		err = fmt.Errorf("unexpected sub-program for command %v", name)
	}
	return
}

func PrintUsage(name string, out io.Writer) (err error) {
	cli, ok := clis[name]
	if !ok {
		return fmt.Errorf("command not found: %v", name)
	}
	err = cli.PrintUsage(out)
	return
}
