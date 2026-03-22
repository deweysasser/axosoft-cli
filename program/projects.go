package program

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// ProjectsCmd lists all Axosoft projects.
type ProjectsCmd struct{}

func (cmd *ProjectsCmd) Run(options *Options) error {
	client, err := newClient(options)
	if err != nil {
		return err
	}

	projects, err := client.ListProjects(context.Background())
	if err != nil {
		return fmt.Errorf("listing projects: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(projects)
}

// ReleasesCmd lists Axosoft releases, optionally filtered by project.
type ReleasesCmd struct {
	ProjectID int `help:"Filter by project ID" default:"0"`
}

func (cmd *ReleasesCmd) Run(options *Options) error {
	client, err := newClient(options)
	if err != nil {
		return err
	}

	releases, err := client.ListReleases(context.Background(), cmd.ProjectID)
	if err != nil {
		return fmt.Errorf("listing releases: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(releases)
}
