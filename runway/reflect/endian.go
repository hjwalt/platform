package reflect

import (
	"encoding/binary"

	"github.com/hjwalt/platform/runway/trusted"
)

func Endian() binary.ByteOrder {
	return trusted.Endian()
}
