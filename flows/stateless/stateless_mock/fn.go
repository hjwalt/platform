package stateless_mock

import (
	"context"
	"errors"
	"strings"

	"github.com/hjwalt/platform/flows/flow"
	"github.com/hjwalt/platform/format"
)

func MockOneToOne(ctx context.Context, m flow.Message[string, string]) (*flow.Message[string, string], error) {
	if strings.ToLower(m.Key) == "mock_error" {
		return nil, ErrMock
	}

	if strings.ToLower(m.Key) == "recoverable" {
		return nil, ErrRecoverable
	}

	if strings.ToLower(m.Key) == "irrecoverable" {
		return nil, ErrIrrecoverable
	}

	if strings.ToLower(m.Key) == "empty" {
		return nil, nil
	}

	if strings.ToLower(m.Key) == "mock_topic" {
		return &flow.Message[string, string]{
			Topic: "mock_topic",
			Key:   m.Key,
			Value: m.Value + "-mock_topic",
		}, nil
	}

	return &flow.Message[string, string]{
		Key:   m.Key,
		Value: m.Value + "-updated",
	}, nil
}

var (
	InputTopic       = flow.GenericTopic("input", format.Broken(), format.Broken())
	OutputTopic      = flow.GenericTopic("output", format.Broken(), format.Broken())
	EmptyTopic       = flow.GenericTopic("", format.Broken(), format.Broken())
	ErrMock          = errors.New("mock")
	ErrRecoverable   = errors.New("recoverable")
	ErrIrrecoverable = errors.New("irrrecoverable")
)
