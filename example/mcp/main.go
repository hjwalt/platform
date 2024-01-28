package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hjwalt/platform/agent/mcp/mcp_brave_search_web"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	runServer(fmt.Sprintf("%s:%d", "localhost", 8000))
}

// GetTimeParams defines the parameters for the cityTime tool.
type GetTimeParams struct {
	City string `json:"city" jsonschema:"City to get time for (nyc, sf, or boston)"`
}

func runServer(url string) {
	// Create an MCP server.
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "time-server",
		Version: "1.0.0",
	}, nil)

	// Add MCP-level logging middleware.
	server.AddReceivingMiddleware(createLoggingMiddleware())

	mcp_brave_search_web.Add(server)

	// Create the streamable HTTP handler.
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return server
	}, nil)

	log.Printf("MCP server listening on %s", url)
	// log.Printf("Available tool: cityTime (cities: nyc, sf, boston)")

	// Start the HTTP server.
	if err := http.ListenAndServe(url, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func createLoggingMiddleware() mcp.Middleware {
	return func(next mcp.MethodHandler) mcp.MethodHandler {
		return func(
			ctx context.Context,
			method string,
			req mcp.Request,
		) (mcp.Result, error) {
			start := time.Now()
			sessionID := req.GetSession().ID()

			// Log request details.
			log.Printf("[REQUEST] Session: %s | Method: %s",
				sessionID,
				method)

			// Call the actual handler.
			result, err := next(ctx, method, req)

			// Log response details.
			duration := time.Since(start)

			if err != nil {
				log.Printf("[RESPONSE] Session: %s | Method: %s | Status: ERROR | Duration: %v | Error: %v",
					sessionID,
					method,
					duration,
					err)
			} else {
				log.Printf("[RESPONSE] Session: %s | Method: %s | Status: OK | Duration: %v",
					sessionID,
					method,
					duration)
			}

			return result, err
		}
	}
}
