package tool_mcp

import (
	"context"

	"github.com/hjwalt/platform/agent"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func AddToMcp[REQ any, RES any](server *mcp.Server, tool agent.SyncTool[REQ, RES]) {
	mcp.AddTool(
		server,
		&mcp.Tool{
			Name:         tool.Name(),
			Title:        tool.Name(),
			Description:  tool.Description(),
			InputSchema:  tool.RequestSchema(),
			OutputSchema: tool.ResultSchema(),
		},
		mcpBehaviour(tool),
	)
}

func mcpBehaviour[REQ any, RES any](tool agent.SyncTool[REQ, RES]) mcp.ToolHandlerFor[REQ, RES] {
	return func(ctx context.Context, req *mcp.CallToolRequest, params REQ) (*mcp.CallToolResult, RES, error) {
		results, err := tool.Apply(ctx, params)
		return nil, results, err
	}
}
