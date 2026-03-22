package axosoft

import (
	"context"
	"fmt"
	"strconv"
)

// GetWorkflowSteps retrieves available workflow transitions for an item.
func (c *Client) GetWorkflowSteps(ctx context.Context, itemType string, id int) ([]WorkflowStep, error) {
	path := fmt.Sprintf("%s/%d/workflow_steps", itemTypePath(itemType), id)

	var result WorkflowStepsResponse
	if err := c.get(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("getting workflow steps for %s/%d: %w", itemType, id, err)
	}
	return result.Data, nil
}

// TransitionItem moves an item to a new workflow step.
func (c *Client) TransitionItem(ctx context.Context, itemType string, id int, stepID int) error {
	path := fmt.Sprintf("%s/%d/workflow_steps", itemTypePath(itemType), id)

	body := map[string]string{
		"step_id": strconv.Itoa(stepID),
	}

	if err := c.post(ctx, path, body, nil); err != nil {
		return fmt.Errorf("transitioning %s/%d to step %d: %w", itemType, id, stepID, err)
	}
	return nil
}
