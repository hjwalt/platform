package task

import (
	"github.com/hjwalt/platform/format"
	"google.golang.org/protobuf/proto"
)

func ProtobufChannel[V proto.Message](channel string) Channel[V] {
	return GenericChannel(channel, format.Protobuf[V]())
}
