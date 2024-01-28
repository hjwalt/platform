package main

import (
	"bufio"
	"errors"
	"os"

	"github.com/hjwalt/platform/model"
	"github.com/hjwalt/platform/store"
	"github.com/hjwalt/platform/write"
	"github.com/hjwalt/runway/configuration"
	"github.com/hjwalt/runway/format"
	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

func main() {
	storageFormat := store.Protojson[*model.ProtobufSchema]()

	// storageFile := OpenFile(name, O_RDWR|O_CREATE|O_TRUNC, 0666)

	conf, err := configuration.Read("model/protobuf.json", storageFormat)
	if err != nil {
		panic(err)
	}

	logger.Info("conf", zap.Any("conf", conf))

	typeMap := Parse(conf)

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
