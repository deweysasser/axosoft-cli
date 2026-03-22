package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deweysasser/axosoft-cli/axosoft"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerItemTools(s *server.MCPServer, client *axosoft.Client) {
	s.AddTool(
		mcp.NewTool("list_items",
			mcp.WithDescription("Search and filter Axosoft work items (features, tasks, incidents, defects)"),
			mcp.WithNumber("project_id", mcp.Description("Filter by project ID")),
			mcp.WithNumber("release_id", mcp.Description("Filter by release ID")),
			mcp.WithString("search", mcp.Description("Search string to filter items")),
			mcp.WithString("type",
				mcp.Description("Item type to list"),
				mcp.Enum("feature", "task", "incident", "defect"),
			),
			mcp.WithNumber("page", mcp.Description("Page number (1-based)")),
			mcp.WithNumber("page_size", mcp.Description("Number of items per page")),
		),
		listItemsHandler(client),
	)

	s.AddTool(
		mcp.NewTool("get_item",
			mcp.WithDescription("Get full details of a single Axosoft work item"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Item ID")),
			mcp.WithString("type", mcp.Required(),
				mcp.Description("Item type"),
				mcp.Enum("feature", "task", "incident", "defect"),
			),
		),
		getItemHandler(client),
	)
}

func listItemsHandler(client *axosoft.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		opts := axosoft.ListItemsOpts{
			ProjectID: int(request.GetFloat("project_id", 0)),
			ReleaseID: int(request.GetFloat("release_id", 0)),
			Search:    request.GetString("search", ""),
			Type:      request.GetString("type", "feature"),
			Page:      int(request.GetFloat("page", 0)),
			PageSize:  int(request.GetFloat("page_size", 0)),
		}

		result, err := client.ListItems(ctx, opts)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshaling result: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}

func getItemHandler(client *axosoft.Client) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireFloat("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		itemType, err := request.RequireString("type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		item, err := client.GetItem(ctx, int(id), itemType)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		data, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("marshaling result: %w", err)
		}

		return mcp.NewToolResultText(string(data)), nil
	}
}
