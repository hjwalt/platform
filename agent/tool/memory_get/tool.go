package memory_get_tool

import (
	"context"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
	"github.com/hjwalt/platform/state"
)

const DefaultName = "memory_get"

type Configuration struct {
	Store state.Store
	Key   string
	Name  string
	Mutex *sync.Mutex
}

type Request struct{}

type Response struct {
	Key     string `json:"key" jsonschema:"storage key for the canonical memory entry"`
	Exists  bool   `json:"exists" jsonschema:"true when memory entry exists and has content"`
	Content string `json:"content" jsonschema:"full memory markdown content"`
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

	t.mutex.Lock()
	defer t.mutex.Unlock()

	s, err := t.store.Read(ctx, t.key)
	if err != nil {
		return Response{}, err
	}

	content := string(s.Value)
	exists := len(s.Value) > 0

	return Response{
		Key:     t.key,
		Exists:  exists,
		Content: content,
	}, nil
}

func (t *tool) Name() string {
	return t.name
}

func (t *tool) Description() string {
	return "get " + t.Name() + " memory entry in markdown format"
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
	return "read current memory markdown content"
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
	if !response.Exists {
		return "memory entry does not exist yet"
	}
	return response.Content
}

func (t *tool) Auto() bool {
	return false
}
