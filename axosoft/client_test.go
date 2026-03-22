package axosoft

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testServer(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := New(srv.URL, "test-token")
	return client, srv
}

func TestListItems(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v5/features", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "test query", r.URL.Query().Get("search_string"))

		resp := ItemsResponse{
			Data: []Item{
				{ID: 1, Name: "Test Feature"},
				{ID: 2, Name: "Another Feature"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	result, err := client.ListItems(context.Background(), ListItemsOpts{
		Search: "test query",
		Type:   "feature",
	})
	require.NoError(t, err)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, "Test Feature", result.Data[0].Name)
}

func TestGetItem(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v5/tasks/42", r.URL.Path)

		resp := struct {
			Data Item `json:"data"`
		}{
			Data: Item{ID: 42, Name: "Test Task"},
		}
		json.NewEncoder(w).Encode(resp)
	})

	item, err := client.GetItem(context.Background(), 42, "task")
	require.NoError(t, err)
	assert.Equal(t, 42, item.ID)
	assert.Equal(t, "Test Task", item.Name)
}

func TestListProjects(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v5/projects", r.URL.Path)

		resp := ProjectsResponse{
			Data: []Project{
				{ID: 1, Name: "Project A"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	projects, err := client.ListProjects(context.Background())
	require.NoError(t, err)
	assert.Len(t, projects, 1)
	assert.Equal(t, "Project A", projects[0].Name)
}

func TestListReleases(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v5/releases", r.URL.Path)
		assert.Equal(t, "5", r.URL.Query().Get("project_id"))

		resp := ReleasesResponse{
			Data: []Release{
				{ID: 10, Name: "v1.0"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	releases, err := client.ListReleases(context.Background(), 5)
	require.NoError(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, "v1.0", releases[0].Name)
}

func TestGetWorkflowSteps(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v5/features/1/workflow_steps", r.URL.Path)

		resp := WorkflowStepsResponse{
			Data: []WorkflowStep{
				{ID: 100, Name: "In Progress"},
				{ID: 101, Name: "Done"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	})

	steps, err := client.GetWorkflowSteps(context.Background(), "feature", 1)
	require.NoError(t, err)
	assert.Len(t, steps, 2)
	assert.Equal(t, "In Progress", steps[0].Name)
}

func TestTransitionItem(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v5/defects/5/workflow_steps", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)

		var body map[string]float64
		json.NewDecoder(r.Body).Decode(&body)
		assert.Equal(t, float64(200), body["step_id"])

		w.WriteHeader(http.StatusOK)
	})

	err := client.TransitionItem(context.Background(), "defect", 5, 200)
	require.NoError(t, err)
}

func TestListItems_InvalidType(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for invalid type")
	})

	_, err := client.ListItems(context.Background(), ListItemsOpts{Type: "sprint"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `unknown item type "sprint"`)
}

func TestGetItem_InvalidType(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for invalid type")
	})

	_, err := client.GetItem(context.Background(), 1, "story")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown item type")
}

func TestGetWorkflowSteps_InvalidType(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for invalid type")
	})

	_, err := client.GetWorkflowSteps(context.Background(), "unknown", 1)
	require.Error(t, err)
}

func TestTransitionItem_InvalidType(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Error("handler should not be called for invalid type")
	})

	err := client.TransitionItem(context.Background(), "epic", 1, 5)
	require.Error(t, err)
}

func TestListReleases_WithoutProjectFilter(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "", r.URL.Query().Get("project_id"), "project_id should not be set")
		json.NewEncoder(w).Encode(ReleasesResponse{Data: []Release{}})
	})

	releases, err := client.ListReleases(context.Background(), 0)
	require.NoError(t, err)
	assert.Empty(t, releases)
}

func TestListItems_AllParams(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, "3", q.Get("project_id"))
		assert.Equal(t, "7", q.Get("release_id"))
		assert.Equal(t, "bug", q.Get("search_string"))
		assert.Equal(t, "2", q.Get("page"))
		assert.Equal(t, "50", q.Get("page_size"))
		json.NewEncoder(w).Encode(ItemsResponse{})
	})

	_, err := client.ListItems(context.Background(), ListItemsOpts{
		Type: "defect", ProjectID: 3, ReleaseID: 7, Search: "bug", Page: 2, PageSize: 50,
	})
	require.NoError(t, err)
}

func TestAPIError_IncludesResponseBody(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error":"access denied"}`))
	})

	_, err := client.ListProjects(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "403")
	assert.Contains(t, err.Error(), "access denied")
}

func TestAPIError(t *testing.T) {
	client, _ := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "invalid token"}`))
	})

	_, err := client.ListProjects(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}
