package axosoft

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListItemsOpts are the options for listing items.
type ListItemsOpts struct {
	ProjectID int
	ReleaseID int
	Search    string
	Type      string // feature, task, incident, defect
	Page      int
	PageSize  int
}

func itemTypePath(itemType string) string {
	switch itemType {
	case "feature":
		return "features"
	case "task":
		return "tasks"
	case "incident":
		return "incidents"
	case "defect":
		return "defects"
	default:
		return "features"
	}
}

// ListItems retrieves work items with optional filtering.
func (c *Client) ListItems(ctx context.Context, opts ListItemsOpts) (*ItemsResponse, error) {
	path := itemTypePath(opts.Type)

	params := url.Values{}
	if opts.ProjectID > 0 {
		params.Set("project_id", strconv.Itoa(opts.ProjectID))
	}
	if opts.ReleaseID > 0 {
		params.Set("release_id", strconv.Itoa(opts.ReleaseID))
	}
	if opts.Search != "" {
		params.Set("search_string", opts.Search)
	}
	if opts.Page > 0 {
		params.Set("page", strconv.Itoa(opts.Page))
	}
	if opts.PageSize > 0 {
		params.Set("page_size", strconv.Itoa(opts.PageSize))
	}

	var result ItemsResponse
	if err := c.get(ctx, path, params, &result); err != nil {
		return nil, fmt.Errorf("listing %s: %w", path, err)
	}
	return &result, nil
}

// GetItem retrieves a single work item by ID and type.
func (c *Client) GetItem(ctx context.Context, id int, itemType string) (*Item, error) {
	path := fmt.Sprintf("%s/%d", itemTypePath(itemType), id)

	var result struct {
		Data Item `json:"data"`
	}
	if err := c.get(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("getting %s/%d: %w", itemType, id, err)
	}
	return &result.Data, nil
}
