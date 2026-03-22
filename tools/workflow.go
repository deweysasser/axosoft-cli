package tools

import (
	"context"
	"fmt"

	"github.com/deweysasser/axosoft-cli/axosoft"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerWorkflowTools(s *server.MCPServer, client AxosoftClient) {
	s.AddTool(
		mcp.NewTool("list_workflow_steps",
			mcp.WithDescription("List available workflow transitions for an Axosoft work item"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Item ID")),
			mcp.WithString("type", mcp.Required(),
				mcp.Description("Item type"),
				mcp.Enum(axosoft.ItemTypes...),
			),
		),
		listWorkflowStepsHandler(client),
	)

	s.AddTool(
		mcp.NewTool("transition_item",
			mcp.WithDescription("Move an Axosoft work item to a new workflow state"),
			mcp.WithNumber("id", mcp.Required(), mcp.Description("Item ID")),
			mcp.WithString("type", mcp.Required(),
				mcp.Description("Item type"),
				mcp.Enum(axosoft.ItemTypes...),
			),
			mcp.WithNumber("step_id", mcp.Required(), mcp.Description("Target workflow step ID")),
		),
		transitionItemHandler(client),
	)
}

func listWorkflowStepsHandler(client AxosoftClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireFloat("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		itemType, err := request.RequireString("type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		steps, err := client.GetWorkflowSteps(ctx, itemType, int(id))
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return jsonToolResult(steps)
	}
}

func transitionItemHandler(client AxosoftClient) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := request.RequireFloat("id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		itemType, err := request.RequireString("type")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}
		stepID, err := request.RequireFloat("step_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if err := client.TransitionItem(ctx, itemType, int(id), int(stepID)); err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Successfully transitioned %s/%d to step %d", itemType, int(id), int(stepID))), nil
	}
}
