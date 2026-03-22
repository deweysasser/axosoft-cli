package tools

import (
	"context"

	"github.com/deweysasser/axosoft-cli/axosoft"
)

// AxosoftClient defines the Axosoft operations used by MCP tool handlers.
type AxosoftClient interface {
	ListItems(ctx context.Context, opts axosoft.ListItemsOpts) (*axosoft.ItemsResponse, error)
	GetItem(ctx context.Context, id int, itemType string) (*axosoft.Item, error)
	GetWorkflowSteps(ctx context.Context, itemType string, id int) ([]axosoft.WorkflowStep, error)
	TransitionItem(ctx context.Context, itemType string, id int, stepID int) error
	ListProjects(ctx context.Context) ([]axosoft.Project, error)
	ListReleases(ctx context.Context, projectID int) ([]axosoft.Release, error)
}
