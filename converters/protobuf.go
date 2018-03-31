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

func prepareMessage(importPaths []string, protoFile string, messageName string) (message proto.Message, err error) {
	parser := protoparse.Parser{
		ImportPaths: importPaths,
	}

	files, err := parser.ParseFiles(protoFile)
	if err != nil {
		return
	}

	descriptor := files[0].FindMessage(messageName)
	message = dynamic.NewMessage(descriptor)
	return
}

type toJSON struct {
	message    proto.Message
	buffer     *bytes.Buffer
	marshaller *jsonpb.Marshaler
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
		err = c.marshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufToJSON(importPaths []string, protoFile, messageName string, enumsAsInts bool, emitDefaults bool, indent string, origName bool) (converter streamconv.Converter, err error) {
	message, err := prepareMessage(importPaths, protoFile, messageName)
	if err != nil {
		return
	}

	marshaller := &jsonpb.Marshaler{
		EnumsAsInts:  enumsAsInts,
		EmitDefaults: emitDefaults,
		Indent:       indent,
		OrigName:     origName,
	}

	converter = &toJSON{message, &bytes.Buffer{}, marshaller}
	return
}

type fromJSON struct {
	message      proto.Message
	unmarshaller *jsonpb.Unmarshaler
	buffer       *proto.Buffer
}

func (c *fromJSON) Convert(src io.Reader) (dst io.Reader, err error) {
	err = c.unmarshaller.Unmarshal(src, c.message)
	if err != nil {
		return
	}

	c.buffer.Reset()
	err = c.buffer.Marshal(c.message)
	return bytes.NewReader(c.buffer.Bytes()), err
}

func NewProtobufFromJSON(importPaths []string, protoFile, messageName string, allowUnknownFields bool) (converter streamconv.Converter, err error) {
	message, err := prepareMessage(importPaths, protoFile, messageName)
	if err != nil {
		return
	}

	unmarshaller := &jsonpb.Unmarshaler{
		AllowUnknownFields: allowUnknownFields,
	}

	converter = &fromJSON{message, unmarshaller, &proto.Buffer{}}
	return
}

type toText struct {
	message    proto.Message
	buffer     *bytes.Buffer
	marshaller *proto.TextMarshaler
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
		err = c.marshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufToText(importPaths []string, protoFile, messageName string, compact bool, expandAny bool) (converter streamconv.Converter, err error) {
	message, err := prepareMessage(importPaths, protoFile, messageName)
	if err != nil {
		return
	}

	marshaller := &proto.TextMarshaler{
		Compact:   compact,
		ExpandAny: expandAny,
	}

	converter = &toText{message, &bytes.Buffer{}, marshaller}
	return
}

type fromText struct {
	message    proto.Message
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

func NewProtobufFromText(importPaths []string, protoFile, messageName string) (converter streamconv.Converter, err error) {
	message, err := prepareMessage(importPaths, protoFile, messageName)
	if err != nil {
		return
	}

	converter = &fromText{message, &bytes.Buffer{}, &proto.Buffer{}}
	return
}

type jsonToText struct {
	message      proto.Message
	unmarshaller *jsonpb.Unmarshaler
	marshaller   *proto.TextMarshaler
}

func (c *jsonToText) Convert(src io.Reader) (dst io.Reader, err error) {
	err = c.unmarshaller.Unmarshal(src, c.message)
	if err != nil {
		return
	}

	pr, pw := io.Pipe()

	go func() {
		err = c.marshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufJSONToText(importPaths []string, protoFile, messageName string, allowUnknownFields bool, compact bool, expandAny bool) (converter streamconv.Converter, err error) {
	message, err := prepareMessage(importPaths, protoFile, messageName)
	if err != nil {
		return
	}

	unmarshaller := &jsonpb.Unmarshaler{
		AllowUnknownFields: allowUnknownFields,
	}

	marshaller := &proto.TextMarshaler{
		Compact:   compact,
		ExpandAny: expandAny,
	}

	converter = &jsonToText{message, unmarshaller, marshaller}
	return
}

type textToJSON struct {
	message    proto.Message
	textBuffer *bytes.Buffer
	marshaller *jsonpb.Marshaler
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
		err = c.marshaller.Marshal(pw, c.message)
		pw.CloseWithError(err)
	}()

	return pr, nil
}

func NewProtobufTextToJSON(importPaths []string, protoFile, messageName string, enumsAsInts bool, emitDefaults bool, indent string, origName bool) (converter streamconv.Converter, err error) {
	message, err := prepareMessage(importPaths, protoFile, messageName)
	if err != nil {
		return
	}

	marshaller := &jsonpb.Marshaler{
		EnumsAsInts:  enumsAsInts,
		EmitDefaults: emitDefaults,
		Indent:       indent,
		OrigName:     origName,
	}

	converter = &textToJSON{message, &bytes.Buffer{}, marshaller}
	return
}
