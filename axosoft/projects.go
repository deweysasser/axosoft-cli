package axosoft

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
)

// ListProjects retrieves all projects.
func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	var result ProjectsResponse
	if err := c.get(ctx, "projects", nil, &result); err != nil {
		return nil, fmt.Errorf("listing projects: %w", err)
	}
	return result.Data, nil
}

// ListReleases retrieves releases, optionally filtered by project.
func (c *Client) ListReleases(ctx context.Context, projectID int) ([]Release, error) {
	var params url.Values
	if projectID > 0 {
		params = url.Values{"project_id": {strconv.Itoa(projectID)}}
	}

	var result ReleasesResponse
	if err := c.get(ctx, "releases", params, &result); err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}
	return result.Data, nil
}
