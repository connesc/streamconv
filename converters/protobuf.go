package converters

import (
	"bytes"
	"streamconv"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

// TODO: handle errors

var parser = protoparse.Parser{}
var marshaller = jsonpb.Marshaler{}
var unmarshaller = jsonpb.Unmarshaler{}

type toJSON struct {
	message *dynamic.Message
	buffer  *bytes.Buffer
}

func (c *toJSON) Convert(src []byte) (dst []byte, err error) {
	err = proto.Unmarshal(src, c.message)
	if err != nil {
		return
	}

	c.buffer.Reset()
	c.buffer.Grow(proto.Size(c.message))
	err = marshaller.Marshal(c.buffer, c.message)
	return c.buffer.Bytes(), err
}

func NewProtobufToJSON(protoFile, messageName string) streamconv.Converter {
	// TODO: handle errors
	files, _ := parser.ParseFiles(protoFile)
	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	return &toJSON{message, &bytes.Buffer{}}
}

type fromJSON struct {
	message *dynamic.Message
	buffer  *proto.Buffer
}

func (c *fromJSON) Convert(src []byte) (dst []byte, err error) {
	err = unmarshaller.Unmarshal(bytes.NewReader(src), c.message)
	if err != nil {
		return
	}

	c.buffer.Reset()
	err = c.buffer.Marshal(c.message)
	return c.buffer.Bytes(), err
}

func NewProtobufFromJSON(protoFile, messageName string) streamconv.Converter {
	// TODO: handle errors
	files, _ := parser.ParseFiles(protoFile)
	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	return &fromJSON{message, &proto.Buffer{}}
}
