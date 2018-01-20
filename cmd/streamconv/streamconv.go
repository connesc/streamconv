package main

import (
	"io"
	"log"
	"os"

	"streamconv"
	"streamconv/converters"
	"streamconv/joiners"
	"streamconv/splitters"
)

func streamConv(splitter streamconv.Splitter, converters []streamconv.Converter, joiner streamconv.Joiner) (err error) {
	for {
		item, err := splitter.ReadItem()
		if err != nil {
			return err
		}

		for _, converter := range converters {
			item, err = converter.Convert(item)
			if err != nil {
				if err == io.EOF {
					err = io.ErrUnexpectedEOF
				}
				return err
			}
		}

		err = joiner.WriteItem(item)
		if err != nil {
			return err
		}
	}
}

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

	// TODO: main CLI

	commands, err := parse(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	if len(commands) < 2 {
		log.Fatalln("not enough commands")
	}

	splitter, err := streamconv.GetSplitter(commands[0], os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}

	converters := make([]streamconv.Converter, len(commands)-2)
	for index, command := range commands[1 : len(commands)-1] {
		converters[index], err = streamconv.GetConverter(command)
		if err != nil {
			log.Fatalln(err)
		}
	}

	joiner, err := streamconv.GetJoiner(commands[len(commands)-1], os.Stdout)
	if err != nil {
		log.Fatalln(err)
	}

	err = streamConv(splitter, converters, joiner)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
