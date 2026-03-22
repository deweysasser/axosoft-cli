package tools

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/deweysasser/axosoft-cli/axosoft"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockClient is a controllable stub of AxosoftClient.
type mockClient struct {
	listItemsResult  *axosoft.ItemsResponse
	getItemResult    *axosoft.Item
	workflowSteps    []axosoft.WorkflowStep
	projects         []axosoft.Project
	releases         []axosoft.Release
	err              error
	transitionCalled bool
	lastTransitionID int
	lastStepID       int
}

func (m *mockClient) ListItems(_ context.Context, _ axosoft.ListItemsOpts) (*axosoft.ItemsResponse, error) {
	return m.listItemsResult, m.err
}
func (m *mockClient) GetItem(_ context.Context, _ int, _ string) (*axosoft.Item, error) {
	return m.getItemResult, m.err
}
func (m *mockClient) GetWorkflowSteps(_ context.Context, _ string, _ int) ([]axosoft.WorkflowStep, error) {
	return m.workflowSteps, m.err
}
func (m *mockClient) TransitionItem(_ context.Context, _ string, id int, stepID int) error {
	m.transitionCalled = true
	m.lastTransitionID = id
	m.lastStepID = stepID
	return m.err
}
func (m *mockClient) ListProjects(_ context.Context) ([]axosoft.Project, error) {
	return m.projects, m.err
}
func (m *mockClient) ListReleases(_ context.Context, _ int) ([]axosoft.Release, error) {
	return m.releases, m.err
}

func callTool(params map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{Arguments: params},
	}
}

func textContent(result *mcp.CallToolResult) string {
	if len(result.Content) == 0 {
		return ""
	}
	tc, ok := result.Content[0].(mcp.TextContent)
	if !ok {
		return ""
	}
	return tc.Text
}

// --- listItemsHandler ---

func TestListItemsHandler_Success(t *testing.T) {
	client := &mockClient{
		listItemsResult: &axosoft.ItemsResponse{
			Data: []axosoft.Item{{ID: 1, Name: "My Feature"}},
		},
	}
	handler := listItemsHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"type": "feature", "search": "my"}))

	require.NoError(t, err)
	assert.False(t, result.IsError)
	var resp axosoft.ItemsResponse
	require.NoError(t, json.Unmarshal([]byte(textContent(result)), &resp))
	assert.Equal(t, 1, resp.Data[0].ID)
}

func TestListItemsHandler_ClientError(t *testing.T) {
	client := &mockClient{err: errors.New("network failure")}
	handler := listItemsHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, textContent(result), "network failure")
}

func TestListItemsHandler_DefaultsTypeToFeature(t *testing.T) {
	client := &mockClient{listItemsResult: &axosoft.ItemsResponse{}}
	handler := listItemsHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
}

// --- getItemHandler ---

func TestGetItemHandler_Success(t *testing.T) {
	client := &mockClient{getItemResult: &axosoft.Item{ID: 42, Name: "Bug"}}
	handler := getItemHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(42), "type": "defect"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	assert.Contains(t, textContent(result), "Bug")
}

func TestGetItemHandler_MissingID(t *testing.T) {
	client := &mockClient{}
	handler := getItemHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"type": "task"}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestGetItemHandler_MissingType(t *testing.T) {
	client := &mockClient{}
	handler := getItemHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(1)}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, textContent(result), "type")
}

func TestGetItemHandler_ClientError(t *testing.T) {
	client := &mockClient{err: errors.New("not found")}
	handler := getItemHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(99), "type": "feature"}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, textContent(result), "not found")
}

// --- listWorkflowStepsHandler ---

func TestListWorkflowStepsHandler_Success(t *testing.T) {
	client := &mockClient{
		workflowSteps: []axosoft.WorkflowStep{{ID: 10, Name: "In Progress"}},
	}
	handler := listWorkflowStepsHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(1), "type": "feature"}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	assert.Contains(t, textContent(result), "In Progress")
}

func TestListWorkflowStepsHandler_MissingID(t *testing.T) {
	handler := listWorkflowStepsHandler(&mockClient{})
	result, err := handler(context.Background(), callTool(map[string]any{"type": "feature"}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestListWorkflowStepsHandler_ClientError(t *testing.T) {
	client := &mockClient{err: errors.New("timeout")}
	handler := listWorkflowStepsHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(1), "type": "feature"}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

// --- transitionItemHandler ---

func TestTransitionItemHandler_Success(t *testing.T) {
	client := &mockClient{}
	handler := transitionItemHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(5), "type": "task", "step_id": float64(200)}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	assert.True(t, client.transitionCalled)
	assert.Equal(t, 5, client.lastTransitionID)
	assert.Equal(t, 200, client.lastStepID)
	assert.Contains(t, textContent(result), "task/5")
	assert.Contains(t, textContent(result), "200")
}

func TestTransitionItemHandler_MissingStepID(t *testing.T) {
	client := &mockClient{}
	handler := transitionItemHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(5), "type": "task"}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.False(t, client.transitionCalled)
}

func TestTransitionItemHandler_ClientError(t *testing.T) {
	client := &mockClient{err: errors.New("forbidden")}
	handler := transitionItemHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"id": float64(5), "type": "task", "step_id": float64(200)}))
	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, textContent(result), "forbidden")
}

// --- listProjectsHandler ---

func TestListProjectsHandler_Success(t *testing.T) {
	client := &mockClient{projects: []axosoft.Project{{ID: 1, Name: "Alpha"}}}
	handler := listProjectsHandler(client)
	result, err := handler(context.Background(), callTool(nil))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	assert.Contains(t, textContent(result), "Alpha")
}

func TestListProjectsHandler_Error(t *testing.T) {
	client := &mockClient{err: errors.New("unauthorized")}
	handler := listProjectsHandler(client)
	result, err := handler(context.Background(), callTool(nil))
	require.NoError(t, err)
	assert.True(t, result.IsError)
}

// --- listReleasesHandler ---

func TestListReleasesHandler_WithProjectFilter(t *testing.T) {
	client := &mockClient{releases: []axosoft.Release{{ID: 10, Name: "v1.0"}}}
	handler := listReleasesHandler(client)
	result, err := handler(context.Background(), callTool(map[string]any{"project_id": float64(5)}))
	require.NoError(t, err)
	assert.False(t, result.IsError)
	assert.Contains(t, textContent(result), "v1.0")
}

func TestListReleasesHandler_NoFilter(t *testing.T) {
	client := &mockClient{releases: []axosoft.Release{}}
	handler := listReleasesHandler(client)
	result, err := handler(context.Background(), callTool(nil))
	require.NoError(t, err)
	assert.False(t, result.IsError)
}

func TestListReleasesHandler_Error(t *testing.T) {
	client := &mockClient{err: errors.New("server unavailable")}
	handler := listReleasesHandler(client)
	result, err := handler(context.Background(), callTool(nil))
	require.NoError(t, err)
	assert.True(t, result.IsError)
	assert.Contains(t, textContent(result), "server unavailable")
}
