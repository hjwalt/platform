package main

import (
	"bufio"
	"errors"
	"os"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/write"
	"github.com/hjwalt/runway/configuration"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

func main() {
	storageFormat := format.Protojson[*model.ProtobufSchema]()

	conf, err := configuration.Read("model/protobuf.json", storageFormat)
	if err != nil {
		panic(err)
	}

	logger.Info("conf", zap.Any("conf", conf))

	typeMap := Parse(conf)

	typeMap["ProtobufMessageField"] = ProtobufMessage(
		"ProtobufMessageField",
		[]*model.ProtobufMessageField{
			ProtobufMessageOneofField(
				"field",
				[]*model.ProtobufMessageField{
					ProtobufMessageBasicField("ProtobufMessageBasicField", "basic_field", 1),
					ProtobufMessageBasicField("ProtobufMessageRepeatedField", "repeated_field", 2),
					ProtobufMessageBasicField("ProtobufMessageMapField", "map_field", 3),
					ProtobufMessageBasicField("ProtobufMessageOneofField", "oneof_field", 4),
				},
			),
		},
	)

	Flatten(conf, typeMap)

	Write("model/protobuf.json", storageFormat, conf)

	protoFile, _ := os.Create("model/protobuf.proto")
	defer protoFile.Close()

	w := bufio.NewWriter(protoFile)
	defer w.Flush()

	write.WriteProtobufSchema(w, conf)

}

func Write[T any](file string, f format.Format[T], val T) error {
	bytes, marshalErr := f.Marshal(val)
	if marshalErr != nil {
		return errors.Join(ErrWriteMarshal, marshalErr)
	}

	writeErr := os.WriteFile(file, bytes, os.ModeExclusive)
	if writeErr != nil {
		return errors.Join(ErrWriteFail, writeErr)
	}

	return nil
}

var ErrWriteMarshal = errors.New("cannot marshal value")
var ErrWriteFail = errors.New("cannot write file")
