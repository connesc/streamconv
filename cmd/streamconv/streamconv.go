package main

import (
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"

	"github.com/connesc/streamconv"
	"github.com/connesc/streamconv/app"
	"github.com/connesc/streamconv/converters"
	"github.com/connesc/streamconv/joiners"
	"github.com/connesc/streamconv/splitters"
)

func main() {
	splitters.RegisterJSONSplitter("json")
	splitters.RegisterSingleSplitter("single")
	splitters.RegisterSimpleSplitter("split")
	splitters.RegisterVarintSplitter("varint")
	splitters.RegisterWindowSplitter("window")
	converters.RegisterBase64Encoder("base64.encode")
	converters.RegisterBase64Decoder("base64.decode")
	converters.RegisterExecutorCommand("exec")
	converters.RegisterProtobufToJSON("protobuf.tojson")
	converters.RegisterProtobufFromJSON("protobuf.fromjson")
	joiners.RegisterSimpleJoiner("join")
	joiners.RegisterVarintJoiner("varint")

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

	if pflag.NArg() != 1 {
		log.Fatalf("invalid number of arguments (expected 1, got %v)\n", pflag.NArg())
	}

	app, err := app.New(pflag.Arg(0))
	if err != nil {
		log.Fatalln(err)
	}

	err = app.Run(os.Stdout, os.Stdin)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
