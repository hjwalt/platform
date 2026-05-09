package shell_tool

import (
	"context"
	"os/exec"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/hjwalt/platform/agent"
	"github.com/hjwalt/platform/agent/llm"
	"github.com/hjwalt/platform/format"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/openai/openai-go/v3"
)

type Configuration struct {
	BaseDir string
}

type Mcp interface {
	agent.Tool
	Behaviour(ctx context.Context, req *mcp.CallToolRequest, params Request) (*mcp.CallToolResult, Response, error)
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

func (t *Tool) Behaviour(ctx context.Context, req *mcp.CallToolRequest, params Request) (*mcp.CallToolResult, Response, error) {
	results, err := t.internal(params)
	return nil, results, err
}

func (t *Tool) internal(params Request) (Response, error) {
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

func (t *Tool) Schema() openai.ChatCompletionToolUnionParam {
	return llm.OpenAiToolSchema[Request](t.Name(), t.Description())
}

func (t *Tool) Execute(input string) (string, error) {
	request, requestParseErr := RequestFormat.Unmarshal([]byte(input))
	if requestParseErr != nil {
		return "", requestParseErr
	}

	response, internalErr := t.internal(request)
	if internalErr != nil {
		return "", internalErr
	}

	return response.Result, nil
}

func (t *Tool) Request(input string) (string, error) {
	request, requestParseErr := RequestFormat.Unmarshal([]byte(input))
	if requestParseErr != nil {
		return "", requestParseErr
	}

	outputBuilder := strings.Builder{}

	outputBuilder.WriteString("execute shell command `")
	outputBuilder.WriteString(request.Command)
	outputBuilder.WriteString("`")

	return outputBuilder.String(), nil
}

func (t *Tool) Auto() bool {
	return false
}

func Add(server *mcp.Server) {
	opts := &jsonschema.ForOptions{}

	in, _ := jsonschema.For[Request](opts)
	out, _ := jsonschema.For[Response](opts)

	instance := Instance(Configuration{
		BaseDir: "/home/hjwalt/Projects/platform/tmp/cmd/",
	})

	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:         instance.Name(),
			Title:        instance.Name(),
			Description:  instance.Description(),
			InputSchema:  in,
			OutputSchema: out,
		},
		instance.Behaviour,
	)
}

func Instance(config Configuration) Mcp {
	return &Tool{
		BaseDir: config.BaseDir,
	}
}

var (
	RequestFormat  = format.Json[Request]()
	ResponseFormat = format.Json[Response]()
)
