package program

import (
	"context"
	"fmt"
)

// ProjectsCmd lists all Axosoft projects.
type ProjectsCmd struct{}

func (cmd *ProjectsCmd) Run(ctx context.Context, options *Options) error {
	client, err := newClient(options)
	if err != nil {
		return err
	}

	projects, err := client.ListProjects(ctx)
	if err != nil {
		return fmt.Errorf("listing projects: %w", err)
	}

	return printJSON(projects)
}

// ReleasesCmd lists Axosoft releases, optionally filtered by project.
type ReleasesCmd struct {
	ProjectID int `help:"Filter by project ID" default:"0"`
}

func (cmd *ReleasesCmd) Run(ctx context.Context, options *Options) error {
	client, err := newClient(options)
	if err != nil {
		return err
	}

	releases, err := client.ListReleases(ctx, cmd.ProjectID)
	if err != nil {
		return fmt.Errorf("listing releases: %w", err)
	}

	return printJSON(releases)
}
