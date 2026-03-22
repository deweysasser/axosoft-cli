package axosoft

import (
	"context"
	"fmt"
)

// GetWorkflowSteps retrieves available workflow transitions for an item.
func (c *Client) GetWorkflowSteps(ctx context.Context, itemType string, id int) ([]WorkflowStep, error) {
	typePath, err := itemTypePath(itemType)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%d/workflow_steps", typePath, id)

	var result WorkflowStepsResponse
	if err := c.get(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("getting workflow steps for %s/%d: %w", itemType, id, err)
	}
	return result.Data, nil
}

// TransitionItem moves an item to a new workflow step.
func (c *Client) TransitionItem(ctx context.Context, itemType string, id int, stepID int) error {
	typePath, err := itemTypePath(itemType)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%d/workflow_steps", typePath, id)

	body := map[string]int{
		"step_id": stepID,
	}

	if err := c.post(ctx, path, body, nil); err != nil {
		return fmt.Errorf("transitioning %s/%d to step %d: %w", itemType, id, stepID, err)
	}
	return nil
}
