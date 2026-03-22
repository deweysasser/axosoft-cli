package auth

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveAndLoad(t *testing.T) {
	// Use a temp dir for XDG_CONFIG_HOME
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	token := &Token{
		AccessToken: "test-access-token",
		ExpiresAt:   time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC),
		BaseURL:     "https://example.axosoft.com",
	}

	err := Save(token)
	require.NoError(t, err)

	// Verify file exists with correct permissions
	path := filepath.Join(tmpDir, "axosoft-cli", "token.json")
	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0600), info.Mode().Perm())

	// Load it back
	loaded, err := Load("https://example.axosoft.com")
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "test-access-token", loaded.AccessToken)
	assert.Equal(t, "https://example.axosoft.com", loaded.BaseURL)
}

func TestLoadWrongURL(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	token := &Token{
		AccessToken: "test-token",
		BaseURL:     "https://example.axosoft.com",
	}
	require.NoError(t, Save(token))

	loaded, err := Load("https://other.axosoft.com")
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

func TestLoadMissing(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	loaded, err := Load("https://example.axosoft.com")
	require.NoError(t, err)
	assert.Nil(t, loaded)
}
