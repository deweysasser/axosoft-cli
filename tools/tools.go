package tools

import (
	"github.com/deweysasser/axosoft-cli/axosoft"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAll registers all Axosoft MCP tools on the server.
func RegisterAll(s *server.MCPServer, client *axosoft.Client) {
	registerItemTools(s, client)
	registerWorkflowTools(s, client)
	registerProjectTools(s, client)
}
