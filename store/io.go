package store

import (
	"errors"
	"os"

	"github.com/hjwalt/runway/format"
)

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

func Read[T any](file string, f format.Format[T]) (T, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return f.Default(), errors.Join(ErrReadFail, err)
	}
	return f.Unmarshal(bytes)
}

var ErrReadFail = errors.New("cannot read file")
