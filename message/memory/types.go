package memory

import (
	"github.com/hjwalt/platform/message"
)

type MemoryMetadata struct {
	Headers map[string]string
}

type MemoryConfiguration struct {
	Channel chan message.Message[MemoryMetadata]
}
