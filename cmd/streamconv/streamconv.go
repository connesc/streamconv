package main

import (
	"io"
	"log"
	"os"

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
	// splitter := readers.NewSplitReader(os.Stdin)
	splitter := readers.NewJSONReader(os.Stdin)
	// splitter := readers.NewSingleReader(os.Stdin)
	// splitter := readers.NewVarintReader(os.Stdin)
	converters := []streamconv.Converter{
		converters.NewProtobufFromJSON("test.proto", "main.SearchRequest"),
		converters.NewBase64Encode(),
		converters.NewBase64Decode(),
		converters.NewProtobufToJSON("test.proto", "main.SearchRequest"),
	}
	joiner := writers.NewJoinWriter(os.Stdout, "\n")
	// joiner := writers.NewVarintWriter(os.Stdout)

	err := streamConv(splitter, converters, joiner)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
