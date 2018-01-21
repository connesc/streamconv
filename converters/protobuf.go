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
