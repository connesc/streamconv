package main

import (
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/connesc/streamconv"
	"github.com/connesc/streamconv/app"
	"github.com/connesc/streamconv/combiners"
	"github.com/connesc/streamconv/converters"
	"github.com/connesc/streamconv/extractors"
	"github.com/connesc/streamconv/transformers"
)

func main() {
	extractors.RegisterJSONExtractor("split.json")
	extractors.RegisterSingleExtractor("single")
	extractors.RegisterSplitExtractor("split")
	extractors.RegisterVarintExtractor("split.varint")
	extractors.RegisterWindowExtractor("window.bytes")
	converters.RegisterBase64Encoder("base64.encode")
	converters.RegisterBase64Decoder("base64.decode")
	converters.RegisterExecutor("exec")
	converters.RegisterJSONIndenter("json.indent")
	converters.RegisterJSONCompactor("json.compact")
	converters.RegisterProtobufToJSON("protobuf.tojson")
	converters.RegisterProtobufFromJSON("protobuf.fromjson")
	converters.RegisterProtobufToText("protobuf.totext")
	converters.RegisterProtobufFromText("protobuf.fromtext")
	converters.RegisterProtobufJSONToText("protobuf.jsontotext")
	converters.RegisterProtobufTextToJSON("protobuf.texttojson")
	transformers.RegisterWindowTransformer("window")
	combiners.RegisterJoinCombiner("join")
	combiners.RegisterCatCombiner("cat")
	combiners.RegisterVarintCombiner("join.varint")

	help := false
	pflag.BoolVarP(&help, "help", "h", help, "print the general help, or the help of the given command")

	pflag.Parse()

	if help {
		if pflag.NArg() == 0 {
			pflag.Usage()
		} else {
			err := streamconv.PrintUsage(pflag.Arg(0), os.Stdout)
			if err != nil {
				log.Fatalln(err)
			}
		}
		return
	}

	if pflag.NArg() > 1 {
		log.Fatalf("too many arguments (expected up to 1, got %v)\n", pflag.NArg())
	}

	program := ""
	if pflag.NArg() == 1 {
		program = pflag.Arg(0)
	}

	app, err := app.Parse(program)
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Run(os.Stdout, os.Stdin)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
