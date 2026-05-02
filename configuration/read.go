package configuration

import (
	"errors"
	"os"

	"github.com/hjwalt/platform/format"
)

var (
	Format = format.Yaml[Configuration]()
)

func Read(file string) (Configuration, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return Format.Default(), errors.Join(ErrReadFail, err)
	}
	return Format.Unmarshal(bytes)
}

var ErrReadFail = errors.New("cannot read file")
