package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// loadAll reads the token store from disk.
func loadAll() (map[string]Token, error) {
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

	// Try new format: map of URL -> token
	var store map[string]Token
	if err := json.Unmarshal(data, &store); err == nil && len(store) > 0 {
		return store, nil
	}

	// Fall back to legacy single-token format for migration
	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("parsing token file: %w", err)
	}
	if token.BaseURL != "" {
		return map[string]Token{token.BaseURL: token}, nil
	}
	return nil, nil
}

// Load reads a cached token for the given base URL.
// Returns nil if no token is found for that URL.
func Load(baseURL string) (*Token, error) {
	baseURL = strings.TrimRight(baseURL, "/")

	store, err := loadAll()
	if err != nil {
		return nil, err
	}

	token, ok := store[baseURL]
	if !ok {
		return nil, nil
	}

	if !token.ExpiresAt.IsZero() && time.Now().After(token.ExpiresAt) {
		return nil, fmt.Errorf("token for %s has expired — run 'axosoft-cli login' again", baseURL)
	}

	return &token, nil
}

// Save writes the token to the config directory, preserving tokens for other URLs.
func Save(token *Token) error {
	token.BaseURL = strings.TrimRight(token.BaseURL, "/")

	store, err := loadAll()
	if err != nil {
		// If the file is corrupt, start fresh
		store = nil
	}
	if store == nil {
		store = make(map[string]Token)
	}
	store[token.BaseURL] = *token

	path, err := tokenPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}

	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling token store: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing token file: %w", err)
	}

	return nil
}
