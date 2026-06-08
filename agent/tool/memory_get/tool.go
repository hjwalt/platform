package memory_get_tool

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
)

const (
	DefaultName    = "memory_get"
	memoryFileName = "memory.md"
)

type Configuration struct {
	RootPath string
	Name     string
	Mutex    *sync.Mutex
}

type Request struct{}

type Response struct {
	Path    string `json:"path" jsonschema:"resolved canonical memory file path"`
	Exists  bool   `json:"exists" jsonschema:"true when canonical memory file already exists"`
	Content string `json:"content" jsonschema:"full memory markdown content"`
}

type tool struct {
	rootPath string
	name     string
	mutex    *sync.Mutex
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
		rootPath: config.RootPath,
		name:     name,
		mutex:    mutex,
	}
}

func (t *tool) Apply(ctx context.Context, request Request) (Response, error) {
	if err := ctx.Err(); err != nil {
		return Response{}, err
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	path := filepath.Join(t.rootPath, memoryFileName)
	bytes, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Response{
			Path:    path,
			Exists:  false,
			Content: "",
		}, nil
	}
	if err != nil {
		return Response{}, err
	}

	return Response{
		Path:    path,
		Exists:  true,
		Content: string(bytes),
	}, nil
}

func (t *tool) Name() string {
	return t.name
}

func (t *tool) Description() string {
	return "Read markdown memory from the canonical memory.md file under the configured root path."
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
		return "memory file does not exist yet"
	}
	return response.Content
}

func (t *tool) Auto() bool {
	return false
}
