package rag

import (
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/runtime"
)

type Store interface {
	runtime.Runtime
	GetAll(id string) ([]agent.Message, error)
	GetFrom(id string, sequence int) ([]agent.Message, error)
	Add(id string, messages []agent.Message) error
	Reset(id string) error
}
