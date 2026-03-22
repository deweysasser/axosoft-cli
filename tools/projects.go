package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deweysasser/axosoft-cli/axosoft"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerProjectTools(s *server.MCPServer, client *axosoft.Client) {
	s.AddTool(
		mcp.NewTool("list_projects",
			mcp.WithDescription("List all Axosoft projects"),
		),
		listProjectsHandler(client),
	)

	s.AddTool(
		mcp.NewTool("list_releases",
			mcp.WithDescription("List Axosoft releases, optionally filtered by project"),
			mcp.WithNumber("project_id", mcp.Description("Filter by project ID")),
		),
		listReleasesHandler(client),
	)
}

func listProjectsHandler(client *axosoft.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projects, err := client.ListProjects(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, err := json.MarshalIndent(projects, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshaling result: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func listReleasesHandler(client *axosoft.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID := int(request.GetFloat("project_id", 0))

		releases, err := client.ListReleases(ctx, projectID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, err := json.MarshalIndent(releases, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshaling result: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}
