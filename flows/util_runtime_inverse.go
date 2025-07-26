package flows

import (
	"github.com/hjwalt/platform/commons/inverse"
)

type Prebuilt interface {
	Register(ci inverse.Container)
}
