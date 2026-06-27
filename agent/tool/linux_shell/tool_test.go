package linux_shell_tool

import (
	"context"
	"testing"

	"github.com/hjwalt/platform/agent"
	harness_container "github.com/hjwalt/platform/agent/util/container"
	"github.com/stretchr/testify/assert"
)

func TestLinuxShellName(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	assert.Equal("linux_shell", tool.Name())
}

func TestLinuxShellDescription(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	assert.NotEmpty(tool.Description())
}

func TestLinuxShellAuto(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	assert.False(tool.Auto())
}

func TestLinuxShellRequestSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	schema := tool.RequestSchema()
	assert.NotNil(schema)
}

func TestLinuxShellResultSchema(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	schema := tool.ResultSchema()
	assert.NotNil(schema)
}

func TestLinuxShellRequestFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	assert.NotNil(tool.RequestFormat())
}

func TestLinuxShellResultFormat(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	assert.NotNil(tool.ResultFormat())
}

func TestLinuxShellDescribeRequest(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	desc := tool.DescribeRequest(Request{Command: "ls -la"})

	assert.Contains(desc, "ls -la")
	assert.Contains(desc, "shell")
}

func TestLinuxShellDescribeResult(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	result := "file1.txt\nfile2.txt"
	desc := tool.DescribeResult(Response{Result: result})

	assert.Equal(result, desc)
}

func TestLinuxShellDescribeResultEmpty(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	desc := tool.DescribeResult(Response{Result: ""})

	assert.Empty(desc)
}

func TestLinuxShellApplyEcho(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	tool := Create(Configuration{BaseDir: "/tmp"})

	resp, err := tool.Apply(ctx, Request{Command: "echo hello world"})

	assert.NoError(err)
	assert.Contains(resp.Result, "hello world")
}

func TestLinuxShellApplyWithBaseDir(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	tool := Create(Configuration{BaseDir: "/tmp"})

	resp, err := tool.Apply(ctx, Request{Command: "pwd"})

	assert.NoError(err)
	assert.Contains(resp.Result, "/tmp")
}

func TestLinuxShellApplyInvalidCommand(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	tool := Create(Configuration{BaseDir: "/tmp"})

	_, err := tool.Apply(ctx, Request{Command: "nonexistentcommand12345"})

	assert.Error(err)
}

func TestLinuxShellApplyEmptyCommand(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()

	tool := Create(Configuration{BaseDir: "/tmp"})

	_, err := tool.Apply(ctx, Request{Command: ""})

	assert.Error(err)
}

func TestLinuxShellCreateReturnsSyncTool(t *testing.T) {
	assert := assert.New(t)

	tool := Create(Configuration{BaseDir: "/tmp"})

	var _ agent.SyncTool[Request, Response] = tool
	assert.NotNil(tool)
}

func TestAddToContainerRegistersLinuxShellTool(t *testing.T) {
	assert := assert.New(t)
	container := harness_container.NewToolContainer()

	AddToContainer(container, Configuration{BaseDir: "/tmp"})

	assert.True(container.Exists(agent.ToolCall{Name: "linux_shell"}))
}
