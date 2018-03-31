package converters

import (
	"bytes"
	"io"

	"github.com/connesc/streamconv"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

var parser = protoparse.Parser{}
var marshaller = jsonpb.Marshaler{}
var unmarshaller = jsonpb.Unmarshaler{}
var textMarshaller = proto.TextMarshaler{}

type toJSON struct {
	message *dynamic.Message
	buffer  *bytes.Buffer
}

func (c *toJSON) Convert(src io.Reader) (dst io.Reader, err error) {
	c.buffer.Reset()
	_, err = c.buffer.ReadFrom(src)
	if err != nil {
		return
	}

	err = proto.Unmarshal(c.buffer.Bytes(), c.message)
	if err != nil {
		return
	}

	pr, pw := io.Pipe()

	go func() {
		err = marshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufToJSON(protoFile, messageName string) (converter streamconv.Converter, err error) {
	files, err := parser.ParseFiles(protoFile)
	if err != nil {
		return
	}

	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	converter = &toJSON{message, &bytes.Buffer{}}
	return
}

type fromJSON struct {
	message *dynamic.Message
	buffer  *proto.Buffer
}

func (c *fromJSON) Convert(src io.Reader) (dst io.Reader, err error) {
	err = unmarshaller.Unmarshal(src, c.message)
	if err != nil {
		return
	}

	c.buffer.Reset()
	err = c.buffer.Marshal(c.message)
	return bytes.NewReader(c.buffer.Bytes()), err
}

func NewProtobufFromJSON(protoFile, messageName string) (converter streamconv.Converter, err error) {
	files, err := parser.ParseFiles(protoFile)
	if err != nil {
		return
	}

	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	converter = &fromJSON{message, &proto.Buffer{}}
	return
}

type toText struct {
	message *dynamic.Message
	buffer  *bytes.Buffer
}

func (c *toText) Convert(src io.Reader) (dst io.Reader, err error) {
	c.buffer.Reset()
	_, err = c.buffer.ReadFrom(src)
	if err != nil {
		return
	}

	err = proto.Unmarshal(c.buffer.Bytes(), c.message)
	if err != nil {
		return
	}

	pr, pw := io.Pipe()

	go func() {
		err = textMarshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufToText(protoFile, messageName string) (converter streamconv.Converter, err error) {
	files, err := parser.ParseFiles(protoFile)
	if err != nil {
		return
	}

	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	converter = &toText{message, &bytes.Buffer{}}
	return
}

type fromText struct {
	message    *dynamic.Message
	textBuffer *bytes.Buffer
	buffer     *proto.Buffer
}

func (c *fromText) Convert(src io.Reader) (dst io.Reader, err error) {
	c.textBuffer.Reset()
	_, err = c.textBuffer.ReadFrom(src)
	if err != nil {
		return
	}

	err = proto.UnmarshalText(c.textBuffer.String(), c.message)
	if err != nil {
		return
	}

	c.buffer.Reset()
	err = c.buffer.Marshal(c.message)
	return bytes.NewReader(c.buffer.Bytes()), err
}

func NewProtobufFromText(protoFile, messageName string) (converter streamconv.Converter, err error) {
	files, err := parser.ParseFiles(protoFile)
	if err != nil {
		return
	}

	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	converter = &fromText{message, &bytes.Buffer{}, &proto.Buffer{}}
	return
}

type jsonToText struct {
	message *dynamic.Message
}

func (c *jsonToText) Convert(src io.Reader) (dst io.Reader, err error) {
	err = unmarshaller.Unmarshal(src, c.message)
	if err != nil {
		return
	}

	pr, pw := io.Pipe()

	go func() {
		err = textMarshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufJSONToText(protoFile, messageName string) (converter streamconv.Converter, err error) {
	files, err := parser.ParseFiles(protoFile)
	if err != nil {
		return
	}

	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	converter = &jsonToText{message}
	return
}

type textToJSON struct {
	message    *dynamic.Message
	textBuffer *bytes.Buffer
}

func (c *textToJSON) Convert(src io.Reader) (dst io.Reader, err error) {
	c.textBuffer.Reset()
	_, err = c.textBuffer.ReadFrom(src)
	if err != nil {
		return
	}

	err = proto.UnmarshalText(c.textBuffer.String(), c.message)
	if err != nil {
		return
	}

	pr, pw := io.Pipe()

	go func() {
		err = marshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufTextToJSON(protoFile, messageName string) (converter streamconv.Converter, err error) {
	files, err := parser.ParseFiles(protoFile)
	if err != nil {
		return
	}

	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	converter = &textToJSON{message, &bytes.Buffer{}}
	return
}
