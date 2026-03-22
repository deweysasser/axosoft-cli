package program

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/deweysasser/axosoft-cli/axosoft"
)

// ItemsCmd lists work items with optional filtering.
type ItemsCmd struct {
	ProjectID int    `help:"Filter by project ID" default:"0"`
	ReleaseID int    `help:"Filter by release ID" default:"0"`
	Search    string `help:"Search string"`
	Type      string `help:"Item type (feature, task, incident, defect)" default:"feature" enum:"feature,task,incident,defect"`
}

func (cmd *ItemsCmd) Run(options *Options) error {
	client, err := newClient(options)
	if err != nil {
		return err
	}

	result, err := client.ListItems(context.Background(), axosoft.ListItemsOpts{
		ProjectID: cmd.ProjectID,
		ReleaseID: cmd.ReleaseID,
		Search:    cmd.Search,
		Type:      cmd.Type,
	})
	if err != nil {
		return fmt.Errorf("listing items: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(result)
}

// ItemCmd retrieves a single work item by ID.
type ItemCmd struct {
	ID   int    `arg:"" help:"Work item ID"`
	Type string `help:"Item type (feature, task, incident, defect)" default:"feature" enum:"feature,task,incident,defect"`
}

func (cmd *ItemCmd) Run(options *Options) error {
	client, err := newClient(options)
	if err != nil {
		return err
	}

	item, err := client.GetItem(context.Background(), cmd.ID, cmd.Type)
	if err != nil {
		return fmt.Errorf("getting item: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(item)
}
