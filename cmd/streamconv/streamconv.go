package main

import (
	"io"
	"log"
	"os"
	"strings"

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

		for _, Converter := range converters {
			item, err = Converter.Convert(item)
			if err != nil {
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
	commands, err := parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	} else {
		for index, command := range commands {
			log.Println(index, ">", strings.Join(command, " "), "<")
		}
	}

	splitter := splitters.NewSimpleSplitter(os.Stdin, "\n")
	// splitter := splitters.NewJSONSplitter(os.Stdin)
	// splitter := splitters.NewSingleSplitter(os.Stdin)
	// splitter := splitters.NewVarintSplitter(os.Stdin)
	// splitter := splitters.NewWindowSplitter(os.Stdin, 3000, 3000, false)
	converters := []streamconv.Converter{
		converters.NewProtobufFromJSON("test.proto", "main.SearchRequest"),
		converters.NewBase64Encode(),
		converters.NewBase64Decode(),
		converters.NewProtobufToJSON("test.proto", "main.SearchRequest"),
	}
	joiner := joiners.NewSimpleJoiner(os.Stdout, "\n")
	// joiner := joiners.NewJoinJoiner(os.Stdout, "")
	// joiner := joiners.NewVarintJoiner(os.Stdout)

	err = streamConv(splitter, converters, joiner)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
