package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerProjectTools(s *server.MCPServer, client AxosoftClient) {
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

func listProjectsHandler(client AxosoftClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projects, err := client.ListProjects(ctx)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return jsonToolResult(projects)
	}
}

func listReleasesHandler(client AxosoftClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		projectID := int(request.GetFloat("project_id", 0))

		releases, err := client.ListReleases(ctx, projectID)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return jsonToolResult(releases)
	}
}
