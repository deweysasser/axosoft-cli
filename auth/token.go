package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Token represents a cached OAuth access token
type Token struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	BaseURL     string    `json:"base_url"`
}

func tokenPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("finding config dir: %w", err)
	}
	return filepath.Join(configDir, "axosoft-cli", "token.json"), nil
}

// Load reads a cached token for the given base URL.
// Returns nil if no token is found or the token is for a different URL.
func Load(baseURL string) (*Token, error) {
	path, err := tokenPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading token file: %w", err)
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("parsing token file: %w", err)
	}

	if token.BaseURL != baseURL {
		return nil, nil
	}

	return &token, nil
}

// Save writes the token to the config directory.
func Save(token *Token) error {
	path, err := tokenPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(token, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling token: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing token file: %w", err)
	}

	return nil
}
