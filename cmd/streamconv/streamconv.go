package main

import (
	"io"
	"log"
	"os"
	"strings"

	"streamconv"
	"streamconv/converters"
	"streamconv/readers"
	"streamconv/writers"
)

func streamConv(splitter streamconv.ItemReader, converters []streamconv.Converter, joiner streamconv.ItemWriter) (err error) {
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

	splitter := readers.NewSplitReader(os.Stdin, "\n")
	// splitter := readers.NewJSONReader(os.Stdin)
	// splitter := readers.NewSingleReader(os.Stdin)
	// splitter := readers.NewVarintReader(os.Stdin)
	// splitter := readers.NewWindowReader(os.Stdin, 3000, 3000, false)
	converters := []streamconv.Converter{
		converters.NewProtobufFromJSON("test.proto", "main.SearchRequest"),
		converters.NewBase64Encode(),
		converters.NewBase64Decode(),
		converters.NewProtobufToJSON("test.proto", "main.SearchRequest"),
	}
	joiner := writers.NewJoinWriter(os.Stdout, "\n")
	// joiner := writers.NewJoinWriter(os.Stdout, "")
	// joiner := writers.NewVarintWriter(os.Stdout)

	err = streamConv(splitter, converters, joiner)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
