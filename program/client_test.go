package program

import (
	"testing"

	"github.com/deweysasser/axosoft-cli/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient_EmptyBaseURL(t *testing.T) {
	opts := &Options{BaseURL: ""}
	_, err := newClient(opts)
	assert.ErrorContains(t, err, "AXOSOFT_URL is required")
}

func TestNewClient_InvalidURLScheme(t *testing.T) {
	opts := &Options{BaseURL: "ftp://example.axosoft.com"}
	_, err := newClient(opts)
	assert.ErrorContains(t, err, "invalid base URL")
}

func TestNewClient_FileURLRejected(t *testing.T) {
	opts := &Options{BaseURL: "file:///etc/passwd"}
	_, err := newClient(opts)
	assert.ErrorContains(t, err, "invalid base URL")
}

func TestNewClient_MissingToken(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	opts := &Options{BaseURL: "https://example.axosoft.com"}
	_, err := newClient(opts)
	assert.ErrorContains(t, err, "run 'axosoft-cli login' first")
}

func TestNewClient_ValidToken(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	token := &auth.Token{
		AccessToken: "valid-token",
		BaseURL:     "https://example.axosoft.com",
	}
	require.NoError(t, auth.Save(token))

	opts := &Options{BaseURL: "https://example.axosoft.com"}
	client, err := newClient(opts)
	require.NoError(t, err)
	assert.NotNil(t, client)
}
