package task

import "github.com/hjwalt/platform/commons/format"

type Channel[V any] interface {
	Name() string
	ValueFormat() format.Format[V]
}
