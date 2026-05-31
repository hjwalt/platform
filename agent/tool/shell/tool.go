package shell_tool

import (
	"context"
	"os/exec"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	tool_mcp "github.com/hjwalt/platform/agent/util/mcp"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/format"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	Name = "linux_shell"
)

type Configuration struct {
	BaseDir string
}

type Request struct {
	Command string `json:"command" jsonschema:"shell command to run with the arguments separated with space"`
}

type Response struct {
	Result string `json:"result" jsonschema:"command output"`
}

type tool struct {
	BaseDir string
}

func (t *tool) Apply(ctx context.Context, params Request) (Response, error) {
	args := strings.Split(params.Command, " ")
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Dir = t.BaseDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Response{}, err
	}

	return Response{
		Result: string(output),
	}, nil
}

func (t *tool) Name() string {
	return Name
}

func (t *tool) Description() string {
	return "Execute command in linux shell with arguments to achieve a specific goal. Assume you have the general tools provided in linux kernels."
}

func (t *tool) RequestFormat() format.Format[Request] {
	return format.Json[Request]()
}

func (t *tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Request](opts)
	return toolSchema
}

func (t *tool) DescribeRequest(request Request) string {
	outputBuilder := strings.Builder{}

	outputBuilder.WriteString("execute shell command `")
	outputBuilder.WriteString(request.Command)
	outputBuilder.WriteString("`")

	return outputBuilder.String()
}

func (t *tool) ResultFormat() format.Format[Response] {
	return format.Json[Response]()
}

func (t *tool) ResultSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Response](opts)
	return toolSchema
}

func (t *tool) DescribeResult(response Response) string {
	return response.Result
}

func (t *tool) Auto() bool {
	return false
}

func Create(config Configuration) agent.SyncTool[Request, Response] {
	return &tool{
		BaseDir: config.BaseDir,
	}
}

func AddToMcp(server *mcp.Server) {
	tool_mcp.AddToMcp(server, Create(Configuration{
		BaseDir: "/home/hjwalt/Projects/platform/tmp/cmd/",
	}))
}

func AddToContainer(container agent.ToolContainer, config Configuration) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(config)))
}
