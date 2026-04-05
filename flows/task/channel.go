package task

import "github.com/hjwalt/platform/format"

type Channel[V any] interface {
	Name() string
	ValueFormat() format.Format[V]
}
