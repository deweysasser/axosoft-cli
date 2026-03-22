package program

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginCmd_MissingBaseURL(t *testing.T) {
	cmd := &LoginCmd{}
	opts := &Options{BaseURL: ""}
	err := cmd.Run(opts)
	assert.ErrorContains(t, err, "AXOSOFT_URL is required")
}

func TestLoginCmd_PasswordFlow_MissingPassword(t *testing.T) {
	cmd := &LoginCmd{Username: "user@example.com", Password: ""}
	opts := &Options{BaseURL: "https://example.axosoft.com", ClientID: "client-id"}
	err := cmd.Run(opts)
	assert.ErrorContains(t, err, "--password is required")
}

func TestLoginCmd_PasswordFlow_MissingClientID(t *testing.T) {
	cmd := &LoginCmd{Username: "user", Password: "pass"}
	opts := &Options{BaseURL: "https://example.axosoft.com", ClientID: ""}
	err := cmd.Run(opts)
	assert.ErrorContains(t, err, "AXOSOFT_CLIENT_ID is required")
}

func TestLoginCmd_PasswordFlow_MissingClientSecret(t *testing.T) {
	cmd := &LoginCmd{Username: "user", Password: "pass", ClientSecret: ""}
	opts := &Options{BaseURL: "https://example.axosoft.com", ClientID: "client-id"}
	err := cmd.Run(opts)
	assert.ErrorContains(t, err, "--client-secret is required")
}

func TestLoginCmd_BrowserFlow_MissingClientID(t *testing.T) {
	cmd := &LoginCmd{}
	opts := &Options{BaseURL: "https://example.axosoft.com", ClientID: ""}
	err := cmd.Run(opts)
	assert.ErrorContains(t, err, "AXOSOFT_CLIENT_ID is required")
}
