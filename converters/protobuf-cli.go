package converters

import (
	"fmt"
	"io"

	"github.com/connesc/streamconv"

	"github.com/spf13/pflag"
)

type protobufToJSONCLI struct {
	name string
}

type protobufToJSONOptions struct {
	importPaths  []string
	enumsAsInts  bool
	emitDefaults bool
	indent       string
	origName     bool
}

func (c *protobufToJSONCLI) newFlagSet() (flags *pflag.FlagSet, options *protobufToJSONOptions) {
	options = &protobufToJSONOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringArrayVarP(&options.importPaths, "proto-path", "I", nil, "directory in which to search for imports (can be repeated, defaults to . if not given)")
	flags.BoolVar(&options.enumsAsInts, "enums-as-ints", false, "render enum values as integers, as opposed to string values")
	flags.BoolVarP(&options.emitDefaults, "defaults", "d", false, "render fields with zero values")
	flags.StringVarP(&options.indent, "indent", "i", "", "a string to indent each level by")
	flags.BoolVarP(&options.origName, "orig-name", "o", false, "use the original (.proto) name for fields")
	return
}

func (c *protobufToJSONCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	if err != nil {
		return
	}

	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *protobufToJSONCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	args = flags.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufToJSON(options.importPaths, args[0], args[1], options.enumsAsInts, options.emitDefaults, options.indent, options.origName)
	}
	return
}

func NewProtobufToJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufToJSONCLI{name}
}

func RegisterProtobufToJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufToJSONCLI(name))
}

type protobufFromJSONCLI struct {
	name string
}

type protobufFromJSONOptions struct {
	importPaths        []string
	allowUnknownFields bool
}

func (c *protobufFromJSONCLI) newFlagSet() (flags *pflag.FlagSet, options *protobufFromJSONOptions) {
	options = &protobufFromJSONOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringArrayVarP(&options.importPaths, "proto-path", "I", nil, "directory in which to search for imports (can be repeated, defaults to . if not given)")
	flags.BoolVarP(&options.allowUnknownFields, "allow-unknown", "k", false, "allow messages to contain unknown fields")
	return
}

func (c *protobufFromJSONCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	if err != nil {
		return
	}

	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *protobufFromJSONCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	args = flags.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufFromJSON(options.importPaths, args[0], args[1], options.allowUnknownFields)
	}
	return
}

func NewProtobufFromJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufFromJSONCLI{name}
}

func RegisterProtobufFromJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufFromJSONCLI(name))
}

type protobufToTextCLI struct {
	name string
}

type protobufToTextOptions struct {
	importPaths []string
	compact     bool
	expandAny   bool
}

func (c *protobufToTextCLI) newFlagSet() (flags *pflag.FlagSet, options *protobufToTextOptions) {
	options = &protobufToTextOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringArrayVarP(&options.importPaths, "proto-path", "I", nil, "directory in which to search for imports (can be repeated, defaults to . if not given)")
	flags.BoolVarP(&options.compact, "compact", "c", false, "use compact text format (one line)")
	flags.BoolVar(&options.expandAny, "expand-any", false, "expand google.protobuf.Any messages of known types")
	return
}

func (c *protobufToTextCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	if err != nil {
		return
	}

	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *protobufToTextCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	args = flags.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufToText(options.importPaths, args[0], args[1], options.compact, options.expandAny)
	}
	return
}

func NewProtobufToTextCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufToTextCLI{name}
}

func RegisterProtobufToText(name string) {
	streamconv.RegisterConverter(name, NewProtobufToTextCLI(name))
}

type protobufFromTextCLI struct {
	name string
}

type protobufFromTextOptions struct {
	importPaths []string
}

func (c *protobufFromTextCLI) newFlagSet() (flags *pflag.FlagSet, options *protobufFromTextOptions) {
	options = &protobufFromTextOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringArrayVarP(&options.importPaths, "proto-path", "I", nil, "directory in which to search for imports (can be repeated, defaults to . if not given)")
	return
}

func (c *protobufFromTextCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	if err != nil {
		return
	}

	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *protobufFromTextCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	args = flags.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufFromText(options.importPaths, args[0], args[1])
	}
	return
}

func NewProtobufFromTextCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufFromTextCLI{name}
}

func RegisterProtobufFromText(name string) {
	streamconv.RegisterConverter(name, NewProtobufFromTextCLI(name))
}

type protobufJSONToTextCLI struct {
	name string
}

type protobufJSONToTextOptions struct {
	importPaths        []string
	allowUnknownFields bool
	compact            bool
	expandAny          bool
}

func (c *protobufJSONToTextCLI) newFlagSet() (flags *pflag.FlagSet, options *protobufJSONToTextOptions) {
	options = &protobufJSONToTextOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringArrayVarP(&options.importPaths, "proto-path", "I", nil, "directory in which to search for imports (can be repeated, defaults to . if not given)")
	flags.BoolVarP(&options.allowUnknownFields, "allow-unknown", "k", false, "allow messages to contain unknown fields")
	flags.BoolVarP(&options.compact, "compact", "c", false, "use compact text format (one line)")
	flags.BoolVar(&options.expandAny, "expand-any", false, "expand google.protobuf.Any messages of known types")
	return
}

func (c *protobufJSONToTextCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	if err != nil {
		return
	}

	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *protobufJSONToTextCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	args = flags.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufJSONToText(options.importPaths, args[0], args[1], options.allowUnknownFields, options.compact, options.expandAny)
	}
	return
}

func NewProtobufJSONToTextCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufJSONToTextCLI{name}
}

func RegisterProtobufJSONToText(name string) {
	streamconv.RegisterConverter(name, NewProtobufJSONToTextCLI(name))
}

type protobufTextToJSONCLI struct {
	name string
}

type protobufTextToJSONOptions struct {
	importPaths  []string
	enumsAsInts  bool
	emitDefaults bool
	indent       string
	origName     bool
}

func (c *protobufTextToJSONCLI) newFlagSet() (flags *pflag.FlagSet, options *protobufTextToJSONOptions) {
	options = &protobufTextToJSONOptions{}
	flags = pflag.NewFlagSet(c.name, pflag.ContinueOnError)
	flags.Usage = func() {}
	flags.StringArrayVarP(&options.importPaths, "proto-path", "I", nil, "directory in which to search for imports (can be repeated, defaults to . if not given)")
	flags.BoolVar(&options.enumsAsInts, "enums-as-ints", false, "render enum values as integers, as opposed to string values")
	flags.BoolVarP(&options.emitDefaults, "defaults", "d", false, "render fields with zero values")
	flags.StringVarP(&options.indent, "indent", "i", "", "a string to indent each level by")
	flags.BoolVarP(&options.origName, "orig-name", "o", false, "use the original (.proto) name for fields")
	return
}

func (c *protobufTextToJSONCLI) PrintUsage(output io.Writer) (err error) {
	_, err = fmt.Fprintln(output, "TODO")
	if err != nil {
		return
	}

	flags, _ := c.newFlagSet()
	flags.SetOutput(output)
	flags.PrintDefaults()
	return
}

func (c *protobufTextToJSONCLI) Parse(args []string) (command streamconv.ConverterCommand, err error) {
	flags, options := c.newFlagSet()
	err = flags.Parse(args)
	if err != nil {
		return
	}

	args = flags.Args()
	if len(args) != 2 {
		return nil, fmt.Errorf("invalid number of arguments (expected 2, got %v)", len(args))
	}

	command = func() (streamconv.Converter, error) {
		return NewProtobufTextToJSON(options.importPaths, args[0], args[1], options.enumsAsInts, options.emitDefaults, options.indent, options.origName)
	}
	return
}

func NewProtobufTextToJSONCLI(name string) (cli streamconv.ConverterCLI) {
	return &protobufTextToJSONCLI{name}
}

func RegisterProtobufTextToJSON(name string) {
	streamconv.RegisterConverter(name, NewProtobufTextToJSONCLI(name))
}
