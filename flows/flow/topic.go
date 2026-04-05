package flow

import "github.com/hjwalt/platform/format"

type Topic[K any, V any] interface {
	Name() string
	KeyFormat() format.Format[K]
	ValueFormat() format.Format[V]
}
