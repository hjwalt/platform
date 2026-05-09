package reflect

import (
	"encoding/binary"

	"github.com/hjwalt/platform/trusted"
)

func Endian() binary.ByteOrder {
	return trusted.Endian()
}
