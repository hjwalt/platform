package shell_tool

import (
	"os/exec"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	agent_tool "github.com/hjwalt/platform/agent/tool"
	tool_string_wrapper "github.com/hjwalt/platform/agent/util/string_wrapper"
	"github.com/hjwalt/platform/format"
	"github.com/modelcontextprotocol/go-sdk/mcp"
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

type Tool struct {
	BaseDir string
}

func (t *Tool) Apply(params Request) (Response, error) {
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

func (t *Tool) Name() string {
	return "linux-shell"
}

func (t *Tool) Description() string {
	return "Execute command in linux shell with arguments to achieve a specific goal. Assume you have the general tools provided in linux kernels."
}

func (t *Tool) RequestFormat() format.Format[Request] {
	return format.Json[Request]()
}

func (t *Tool) RequestSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Request](opts)
	return toolSchema
}

func (t *Tool) DescribeRequest(request Request) string {
	outputBuilder := strings.Builder{}

	outputBuilder.WriteString("execute shell command `")
	outputBuilder.WriteString(request.Command)
	outputBuilder.WriteString("`")

	return outputBuilder.String()
}

func (t *Tool) ResultFormat() format.Format[Response] {
	return format.Json[Response]()
}

func (t *Tool) ResultSchema() *jsonschema.Schema {
	opts := &jsonschema.ForOptions{}
	toolSchema, _ := jsonschema.For[Response](opts)
	return toolSchema
}

func (t *Tool) DescribeResult(response Response) string {
	return response.Result
}

func (t *Tool) Auto() bool {
	return false
}

func Create(config Configuration) agent_tool.Sync[Request, Response] {
	return &Tool{
		BaseDir: config.BaseDir,
	}
}

func AddToMcp(server *mcp.Server) {
	agent_tool.AddToMcp(server, Create(Configuration{
		BaseDir: "/home/hjwalt/Projects/platform/tmp/cmd/",
	}))
}

func AddToContainer(container agent_tool.Container, config Configuration) {
	container.AddSync(tool_string_wrapper.StringWrapSync(Create(config)))
}
