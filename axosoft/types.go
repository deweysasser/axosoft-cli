package axosoft

import "encoding/json"

// itemTypePaths is the single source of truth for item type → API path mapping.
var itemTypePaths = map[string]string{
	"feature":  "features",
	"task":     "tasks",
	"incident": "incidents",
	"defect":   "defects",
}

// ItemTypes returns the canonical list of Axosoft work item types.
func ItemTypes() []string {
	return []string{"feature", "task", "incident", "defect"}
}

// Item represents an Axosoft work item (feature, task, incident, or defect).
type Item struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	ItemType     string         `json:"item_type,omitempty"`
	Number       json.Number    `json:"number,omitempty"`
	Description  string         `json:"description,omitempty"`
	Priority     *NamedResource `json:"priority,omitempty"`
	Status       *NamedResource `json:"status,omitempty"`
	Project      *NamedResource `json:"project,omitempty"`
	Release      *NamedResource `json:"release,omitempty"`
	AssignedTo   *NamedResource `json:"assigned_to,omitempty"`
	Workflow     *NamedResource `json:"workflow_step,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

// NamedResource is a common Axosoft pattern for referenced objects.
type NamedResource struct {
	ID   int    `json:"id"`
	Name string `json:"name,omitempty"`
}

// ItemsResponse is the API response for listing items.
type ItemsResponse struct {
	Data     []Item   `json:"data"`
	Metadata Metadata `json:"metadata,omitempty"`
}

// Metadata contains pagination info.
type Metadata struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}

// Project represents an Axosoft project.
type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ProjectsResponse is the API response for listing projects.
type ProjectsResponse struct {
	Data     []Project `json:"data"`
	Metadata Metadata  `json:"metadata,omitempty"`
}

// Release represents an Axosoft release.
type Release struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
}

// ReleasesResponse is the API response for listing releases.
type ReleasesResponse struct {
	Data     []Release `json:"data"`
	Metadata Metadata  `json:"metadata,omitempty"`
}

// WorkflowStep represents an available workflow transition.
type WorkflowStep struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// WorkflowStepsResponse is the API response for listing workflow steps.
type WorkflowStepsResponse struct {
	Data []WorkflowStep `json:"data"`
}
