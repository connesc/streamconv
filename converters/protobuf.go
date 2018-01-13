package converters

import (
	"streamconv"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

// TODO: handle errors and reuse buffers

var parser = protoparse.Parser{}

type toJSON struct {
	message *dynamic.Message
}

func (c toJSON) Convert(src []byte) (dst []byte, err error) {
	_ = proto.Unmarshal(src, c.message)
	marshaller := jsonpb.Marshaler{}
	json, _ := marshaller.MarshalToString(c.message)

	return []byte(json), nil
}

func NewProtobufToJSON(protoFile, messageName string) streamconv.Converter {
	files, _ := parser.ParseFiles(protoFile)
	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	return toJSON{message}
}

type fromJSON struct {
	message *dynamic.Message
}

func (c fromJSON) Convert(src []byte) (dst []byte, err error) {
	_ = jsonpb.UnmarshalString(string(src), c.message)
	protobuf, _ := proto.Marshal(c.message)
	return protobuf, nil
}

func NewProtobufFromJSON(protoFile, messageName string) streamconv.Converter {
	files, _ := parser.ParseFiles(protoFile)
	descriptor := files[0].FindMessage(messageName)
	message := dynamic.NewMessage(descriptor)

	return fromJSON{message}
}
