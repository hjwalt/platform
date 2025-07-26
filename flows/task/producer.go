package task

import (
	"context"

	"github.com/hjwalt/platform/commons/structure"
)

type Producer interface {
	Produce(context.Context, Message[structure.Bytes]) error
	Start() error
	Stop()
}
