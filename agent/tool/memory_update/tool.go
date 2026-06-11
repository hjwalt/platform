package memory_update_tool

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/state"
)

const DefaultName = "memory_update"

type UpdateMode string

const (
	UpdateModeReplace UpdateMode = "replace"
	UpdateModeAppend  UpdateMode = "append"
)

type Configuration struct {
	Store state.Store
	Key   string
	Name  string
	Mutex *sync.Mutex
}

type Request struct {
	Content string     `json:"content" jsonschema:"memory markdown content to write"`
	Mode    UpdateMode `json:"mode" jsonschema:"write mode: replace or append"`
}

type Response struct {
	Key   string     `json:"key" jsonschema:"storage key for the canonical memory entry"`
	Mode  UpdateMode `json:"mode" jsonschema:"effective write mode used"`
	Bytes int        `json:"bytes" jsonschema:"number of bytes in resulting memory entry"`
}

type tool struct {
	store state.Store
	key   string
	name  string
	mutex *sync.Mutex
}

func Create(config Configuration) agent.SyncTool[Request, Response] {
	name := config.Name
	if name == "" {
		name = DefaultName
	}

	mutex := config.Mutex
	if mutex == nil {
		mutex = &sync.Mutex{}
	}

	return &tool{
		store: config.Store,
		key:   config.Key,
		name:  name,
		mutex: mutex,
	}
}

func (t *tool) Apply(ctx context.Context, request Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	mode, modeErr := resolveUpdateMode(request.Mode)
	if modeErr != nil {
		return Response{}, modeErr
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	content := request.Content
	if mode == UpdateModeAppend {
		existing, readErr := t.store.Read(ctx, t.key)
		if readErr != nil {
			return Response{}, readErr
		}
		content = string(existing.Value) + request.Content
	}

	writeErr := t.store.Write(ctx, state.State{
		Id:        t.key,
		Value:     []byte(content),
		Timestamp: time.Now(),
	})
	if writeErr != nil {
		return Response{}, writeErr
	}

	return Response{
		Key:   t.key,
		Mode:  mode,
		Bytes: len([]byte(content)),
	}, nil
}

func (t *tool) Name() string {
	return t.name
}

func (t *tool) Description() string {
	return "update " + t.Name() + " memory entry in markdown format"
}

func (t *tool) RequestFormat() format.Format[Request] {
	return format.Json[Request]()
}

func (t *tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	schema, _ := jsonschema.For[Request](opts)
	return schema
}

func (t *tool) DescribeRequest(request Request) string {
	mode, modeErr := resolveUpdateMode(request.Mode)
	if modeErr != nil {
		return "write memory with invalid mode"
	}
	return "write memory content in " + string(mode) + " mode"
}

func (t *tool) ResultFormat() format.Format[Response] {
	return format.Json[Response]()
}

func (t *tool) ResultSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	schema, _ := jsonschema.For[Response](opts)
	return schema
}

func (t *tool) DescribeResult(response Response) string {
	return "updated memory entry"
}

func (t *tool) Auto() bool {
	return false
}

func resolveUpdateMode(mode UpdateMode) (UpdateMode, error) {
	if mode == "" {
		return UpdateModeReplace, nil
	}
	if mode == UpdateModeReplace || mode == UpdateModeAppend {
		return mode, nil
	}
	return "", ErrInvalidUpdateMode
}

var ErrInvalidUpdateMode = errors.New("memory update mode must be replace or append")
