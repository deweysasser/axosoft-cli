package tools

import (
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all Axosoft MCP tools on the server.
func RegisterAll(s *server.MCPServer, client AxosoftClient) {
	registerItemTools(s, client)
	registerWorkflowTools(s, client)
	registerProjectTools(s, client)
}
