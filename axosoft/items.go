package axosoft

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

const (
	defaultAutoPageSize = 100
	maxAutoPages        = 100
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

func itemTypePath(itemType string) (string, error) {
	p, ok := itemTypePaths[itemType]
	if !ok {
		return "", fmt.Errorf("unknown item type %q", itemType)
	}
	return p, nil
}

// ListItems retrieves work items with optional filtering.
// If opts.Page > 0, fetches only that page. If opts.Page == 0, auto-paginates
// to retrieve all results.
func (c *Client) ListItems(ctx context.Context, opts ListItemsOpts) (*ItemsResponse, error) {
	if _, err := itemTypePath(opts.Type); err != nil {
		return nil, err
	}

	// Explicit page requested — single-page fetch.
	if opts.Page > 0 {
		return c.listItemsPage(ctx, opts)
	}

	// Auto-paginate: loop until we have all items.
	if opts.PageSize <= 0 {
		opts.PageSize = defaultAutoPageSize
	}

	var allItems []Item
	for page := 1; ; page++ {
		if page > maxAutoPages {
			return nil, fmt.Errorf("auto-pagination safety limit (%d pages) exceeded", maxAutoPages)
		}
		opts.Page = page
		resp, err := c.listItemsPage(ctx, opts)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, resp.Data...)
		if len(resp.Data) == 0 || len(allItems) >= resp.Metadata.TotalCount {
			resp.Data = allItems
			return resp, nil
		}
	}
}

// listItemsPage fetches a single page of items.
func (c *Client) listItemsPage(ctx context.Context, opts ListItemsOpts) (*ItemsResponse, error) {
	path, err := itemTypePath(opts.Type)
	if err != nil {
		return nil, err
	}

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
	typePath, err := itemTypePath(itemType)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("%s/%d", typePath, id)

	var result struct {
		Data Item `json:"data"`
	}
	if err := c.get(ctx, path, nil, &result); err != nil {
		return nil, fmt.Errorf("getting %s/%d: %w", itemType, id, err)
	}
	return &result.Data, nil
}
