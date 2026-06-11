package memory_tool

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/state"
)

const (
	MemoryKey  = "memory"
	MemoryName = "memory"
)

type Configuration struct {
	Key string
}

type Operation string

const (
	OperationGet    Operation = "get"
	OperationUpdate Operation = "update"
	OperationClear  Operation = "clear"
)

type UpdateMode string

const (
	UpdateModeReplace UpdateMode = "replace"
	UpdateModeAppend  UpdateMode = "append"
)

type Request struct {
	Operation Operation  `json:"operation" jsonschema:"required operation: get, update, or clear"`
	Content   string     `json:"content" jsonschema:"memory markdown content to write when operation is update"`
	Mode      UpdateMode `json:"mode" jsonschema:"update mode for operation=update: replace or append"`
}

type Response struct {
	Key       string     `json:"key" jsonschema:"storage key for the canonical memory entry"`
	Operation Operation  `json:"operation" jsonschema:"operation that was executed"`
	Exists    bool       `json:"exists" jsonschema:"true when memory entry exists and has content"`
	Content   string     `json:"content" jsonschema:"full memory markdown content for operation=get"`
	Mode      UpdateMode `json:"mode" jsonschema:"effective write mode used for operation=update"`
	Bytes     int        `json:"bytes" jsonschema:"number of bytes in resulting memory entry for operation=update"`
	Cleared   bool       `json:"cleared" jsonschema:"true when memory content was cleared for operation=clear"`
}

type tool struct {
	store state.Store
	key   string
	name  string
	mutex *sync.Mutex
}

func Create(config Configuration, store state.Store) agent.SyncTool[Request, Response] {
	name := toolName(config.Key, MemoryName)
	key := memoryKey(config.Key)

	return &tool{
		store: store,
		key:   key,
		name:  name,
		mutex: &sync.Mutex{},
	}
}

func (t *tool) Apply(ctx context.Context, request Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}
	if t.store == nil {
		return Response{}, ErrNilStore
	}

	operation, err := resolveOperation(request.Operation)
	if err != nil {
		return Response{}, err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	switch operation {
	case OperationGet:
		return t.applyGet(ctx)
	case OperationUpdate:
		return t.applyUpdate(ctx, request)
	case OperationClear:
		return t.applyClear(ctx)
	default:
		return Response{}, ErrInvalidOperation
	}
}

func (t *tool) Name() string {
	return t.name
}

func (t *tool) Description() string {
	return "manage memory content with one tool using operation=get|update|clear"
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
	switch request.Operation {
	case OperationGet:
		return "read current memory markdown content"
	case OperationUpdate:
		mode, err := resolveUpdateMode(request.Mode)
		if err != nil {
			return "write memory with invalid mode"
		}
		return "write memory content in " + string(mode) + " mode"
	case OperationClear:
		return "clear all memory content"
	default:
		return "run memory operation"
	}
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
	switch response.Operation {
	case OperationGet:
		if !response.Exists {
			return "memory entry does not exist yet"
		}
		return response.Content
	case OperationUpdate:
		return "updated memory entry"
	case OperationClear:
		return "memory content cleared"
	default:
		return "memory operation completed"
	}
}

func (t *tool) Auto() bool {
	return false
}

func (t *tool) applyGet(ctx context.Context) (Response, error) {
	s, err := t.store.Read(ctx, t.key)
	if err != nil {
		return Response{}, err
	}

	content := string(s.Value)
	exists := len(s.Value) > 0

	return Response{
		Key:       t.key,
		Operation: OperationGet,
		Exists:    exists,
		Content:   content,
	}, nil
}

func (t *tool) applyUpdate(ctx context.Context, request Request) (Response, error) {
	mode, err := resolveUpdateMode(request.Mode)
	if err != nil {
		return Response{}, err
	}

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
		Key:       t.key,
		Operation: OperationUpdate,
		Mode:      mode,
		Bytes:     len([]byte(content)),
	}, nil
}

func (t *tool) applyClear(ctx context.Context) (Response, error) {
	if deleteErr := t.store.Delete(ctx, t.key); deleteErr != nil {
		return Response{}, deleteErr
	}

	return Response{
		Key:       t.key,
		Operation: OperationClear,
		Cleared:   true,
	}, nil
}

func AddToContainer(container agent.ToolContainer, config Configuration, store state.Store) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(config, store)))
}

func toolName(prefix string, base string) string {
	if prefix == "" {
		return base
	}
	return prefix + "_" + base
}

func memoryKey(prefix string) string {
	if prefix == "" {
		return MemoryKey
	}
	return prefix
}

var (
	prefixPattern       = regexp.MustCompile(`^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$`)
	ErrNilStore         = errors.New("memory store cannot be nil")
	ErrInvalidPrefix    = errors.New("memory tool prefix must match ^[a-zA-Z0-9]+(?:_[a-zA-Z0-9]+)*$")
	ErrInvalidMode      = errors.New("memory update mode must be replace or append")
	ErrInvalidOperation = errors.New("memory operation must be get, update, or clear")
)

func Validate(config Configuration, store state.Store) error {
	if store == nil {
		return ErrNilStore
	}
	if config.Key == "" {
		return nil
	}
	if !prefixPattern.MatchString(config.Key) {
		return ErrInvalidPrefix
	}
	return nil
}

func resolveUpdateMode(mode UpdateMode) (UpdateMode, error) {
	if mode == "" {
		return UpdateModeReplace, nil
	}
	if mode == UpdateModeReplace || mode == UpdateModeAppend {
		return mode, nil
	}
	return "", ErrInvalidMode
}

func resolveOperation(operation Operation) (Operation, error) {
	normalized := Operation(strings.ToLower(strings.TrimSpace(string(operation))))
	if normalized == OperationGet || normalized == OperationUpdate || normalized == OperationClear {
		return normalized, nil
	}
	return "", ErrInvalidOperation
}
