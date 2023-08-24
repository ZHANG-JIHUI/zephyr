package remote

import (
	"google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type Serializer interface {
	Serialize(msg any) ([]byte, error)
	TypeName(any) string
}

type Deserializer interface {
	Deserialize([]byte, string) (any, error)
}

type (
	ProtoSerializer   struct{}
	ProtoDeserializer struct{}
)

func (ProtoSerializer) Serialize(msg any) ([]byte, error) {
	return proto.Marshal(msg.(proto.Message))
}

func (ProtoSerializer) TypeName(msg any) string {
	return string(proto.MessageName(msg.(proto.Message)))
}

func (ProtoDeserializer) Deserialize(data []byte, name string) (any, error) {
	fullName := protoreflect.FullName(name)
	n, err := protoregistry.GlobalTypes.FindMessageByName(fullName)
	if err != nil {
		return nil, err
	}
	pm := n.New().Interface()
	err = proto.Unmarshal(data, pm)
	return pm, err
}
