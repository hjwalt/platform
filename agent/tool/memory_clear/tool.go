package memory_clear_tool

import (
	"context"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/format"
)

const (
	DefaultName    = "memory_clear"
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
	Cleared bool   `json:"cleared" jsonschema:"true when memory content was cleared"`
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
	if mkdirErr := os.MkdirAll(t.rootPath, 0o755); mkdirErr != nil {
		return Response{}, mkdirErr
	}

	if writeErr := atomicWrite(path, []byte("")); writeErr != nil {
		return Response{}, writeErr
	}

	return Response{
		Path:    path,
		Cleared: true,
	}, nil
}

func (t *tool) Name() string {
	return t.name
}

func (t *tool) Description() string {
	return "Clear markdown memory by truncating canonical memory.md under the configured root path."
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
	return "clear all memory content"
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
	return "memory content cleared"
}

func (t *tool) Auto() bool {
	return false
}

func atomicWrite(path string, content []byte) error {
	tempFile, createErr := os.CreateTemp(filepath.Dir(path), "memory-*.tmp")
	if createErr != nil {
		return createErr
	}

	tempName := tempFile.Name()
	defer os.Remove(tempName)

	if _, writeErr := tempFile.Write(content); writeErr != nil {
		tempFile.Close()
		return writeErr
	}

	if syncErr := tempFile.Sync(); syncErr != nil {
		tempFile.Close()
		return syncErr
	}

	if closeErr := tempFile.Close(); closeErr != nil {
		return closeErr
	}

	if renameErr := os.Rename(tempName, path); renameErr != nil {
		return renameErr
	}

	return nil
}
