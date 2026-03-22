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

func TestMultipleInstances(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	token1 := &Token{
		AccessToken: "token-one",
		BaseURL:     "https://one.axosoft.com",
	}
	token2 := &Token{
		AccessToken: "token-two",
		BaseURL:     "https://two.axosoft.com",
	}

	require.NoError(t, Save(token1))
	require.NoError(t, Save(token2))

	loaded1, err := Load("https://one.axosoft.com")
	require.NoError(t, err)
	require.NotNil(t, loaded1)
	assert.Equal(t, "token-one", loaded1.AccessToken)

	loaded2, err := Load("https://two.axosoft.com")
	require.NoError(t, err)
	require.NotNil(t, loaded2)
	assert.Equal(t, "token-two", loaded2.AccessToken)
}

func TestLoadExpiredToken(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	token := &Token{
		AccessToken: "old-token",
		ExpiresAt:   time.Now().Add(-1 * time.Hour),
		BaseURL:     "https://example.axosoft.com",
	}
	require.NoError(t, Save(token))

	_, err := Load("https://example.axosoft.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expired")
	assert.Contains(t, err.Error(), "axosoft-cli login")
}

func TestSavePreservesOtherURLs(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	require.NoError(t, Save(&Token{AccessToken: "tok-a", BaseURL: "https://a.axosoft.com"}))
	require.NoError(t, Save(&Token{AccessToken: "tok-a-v2", BaseURL: "https://a.axosoft.com"}))
	require.NoError(t, Save(&Token{AccessToken: "tok-b", BaseURL: "https://b.axosoft.com"}))

	loaded, err := Load("https://a.axosoft.com")
	require.NoError(t, err)
	assert.Equal(t, "tok-a-v2", loaded.AccessToken)

	loaded, err = Load("https://b.axosoft.com")
	require.NoError(t, err)
	assert.Equal(t, "tok-b", loaded.AccessToken)
}

func TestLegacyMigration_NoBaseURL(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	dir := filepath.Join(tmpDir, "axosoft-cli")
	require.NoError(t, os.MkdirAll(dir, 0700))
	legacy := `{"access_token":"orphan","expires_at":"2030-01-01T00:00:00Z","base_url":""}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "token.json"), []byte(legacy), 0600))

	loaded, err := Load("https://example.axosoft.com")
	require.NoError(t, err)
	assert.Nil(t, loaded)
}

func TestLegacyMigration(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	// Write a legacy single-token file
	dir := filepath.Join(tmpDir, "axosoft-cli")
	require.NoError(t, os.MkdirAll(dir, 0700))
	legacy := `{"access_token":"old-token","expires_at":"2030-01-01T00:00:00Z","base_url":"https://legacy.axosoft.com"}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "token.json"), []byte(legacy), 0600))

	// Should still load the legacy token
	loaded, err := Load("https://legacy.axosoft.com")
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "old-token", loaded.AccessToken)

	// Saving a new token should migrate to new format and preserve legacy
	newToken := &Token{
		AccessToken: "new-token",
		BaseURL:     "https://new.axosoft.com",
	}
	require.NoError(t, Save(newToken))

	// Both should be loadable
	loaded, err = Load("https://legacy.axosoft.com")
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "old-token", loaded.AccessToken)

	loaded, err = Load("https://new.axosoft.com")
	require.NoError(t, err)
	require.NotNil(t, loaded)
	assert.Equal(t, "new-token", loaded.AccessToken)
}
