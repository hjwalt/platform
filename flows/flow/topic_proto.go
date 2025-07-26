package flow

import (
	"github.com/hjwalt/platform/commons/format"
	"google.golang.org/protobuf/proto"
)

func ProtobufTopic[K proto.Message, V proto.Message](topic string) Topic[K, V] {
	return GenericTopic[K, V](topic, format.Protobuf[K](), format.Protobuf[V]())
}
