package memory_update_tool

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
	DefaultName    = "memory_update"
	memoryFileName = "memory.md"
)

type UpdateMode string

const (
	UpdateModeReplace UpdateMode = "replace"
	UpdateModeAppend  UpdateMode = "append"
)

type Configuration struct {
	BaseDir  string
	FileName string
	Name     string
	Mutex    *sync.Mutex
}

type Request struct {
	Content string     `json:"content" jsonschema:"memory markdown content to write"`
	Mode    UpdateMode `json:"mode" jsonschema:"write mode: replace or append"`
}

type Response struct {
	Path  string     `json:"path" jsonschema:"resolved canonical memory file path"`
	Mode  UpdateMode `json:"mode" jsonschema:"effective write mode used"`
	Bytes int        `json:"bytes" jsonschema:"number of bytes in resulting memory file"`
}

type tool struct {
	baseDir  string
	fileName string
	name     string
	mutex    *sync.Mutex
}

func Create(config Configuration) agent.SyncTool[Request, Response] {
	name := config.Name
	if name == "" {
		name = DefaultName
	}

	fileName := config.FileName
	if fileName == "" {
		fileName = memoryFileName
	}

	mutex := config.Mutex
	if mutex == nil {
		mutex = &sync.Mutex{}
	}

	return &tool{
		baseDir:  config.BaseDir,
		fileName: fileName,
		name:     name,
		mutex:    mutex,
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

	path := filepath.Join(t.baseDir, t.fileName)
	if mkdirErr := os.MkdirAll(t.baseDir, 0o755); mkdirErr != nil {
		return Response{}, mkdirErr
	}

	content := request.Content
	if mode == UpdateModeAppend {
		existing, readErr := os.ReadFile(path)
		if readErr != nil && !errors.Is(readErr, os.ErrNotExist) {
			return Response{}, readErr
		}
		if readErr == nil {
			content = string(existing) + request.Content
		}
	}

	if writeErr := atomicWrite(path, []byte(content)); writeErr != nil {
		return Response{}, writeErr
	}

	return Response{
		Path:  path,
		Mode:  mode,
		Bytes: len([]byte(content)),
	}, nil
}

func (t *tool) Name() string {
	return t.name
}

func (t *tool) Description() string {
	return "Write markdown memory to the configured canonical memory file with deterministic replace or append mode."
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
	return "updated memory file"
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

var ErrInvalidUpdateMode = errors.New("memory update mode must be replace or append")
