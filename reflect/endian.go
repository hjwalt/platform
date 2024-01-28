package reflect

import (
	"encoding/binary"

	"github.com/hjwalt/platform/commons/trusted"
)

func Endian() binary.ByteOrder {
	return trusted.Endian()
}
