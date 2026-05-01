package rag

import (
	"sync"

	"github.com/hjwalt/platform/agent"
)

func Memory() Store {
	return &MemoryStore{
		Messages: make(map[string][]agent.Message),
	}
}

type MemoryStore struct {
	Messages map[string][]agent.Message
	mu       sync.Mutex
}

func (m *MemoryStore) Start() error {
	return nil
}

func (m *MemoryStore) Stop() {

}

func (m *MemoryStore) GetAll(id string) ([]agent.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v, exists := m.Messages[id]; exists {
		return v, nil
	}
	return []agent.Message{}, nil
}

func (m *MemoryStore) GetFrom(id string, sequence int) ([]agent.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if v, exists := m.Messages[id]; exists {
		return v[sequence:], nil
	}
	return []agent.Message{}, nil
}

func (m *MemoryStore) Add(id string, messages []agent.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.Messages[id]; !exists {
		m.Messages[id] = make([]agent.Message, 0)
	}

	m.Messages[id] = append(m.Messages[id], messages...)
	return nil
}

func (m *MemoryStore) Reset(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.Messages[id] = make([]agent.Message, 0)
	return nil
}
