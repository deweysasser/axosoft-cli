package program

import (
	"fmt"
	"net/url"

	"github.com/deweysasser/axosoft-cli/auth"
	"github.com/deweysasser/axosoft-cli/axosoft"
)

// validateBaseURL checks that the URL is non-empty and uses http or https.
func validateBaseURL(baseURL string) error {
	if baseURL == "" {
		return fmt.Errorf("--base-url or AXOSOFT_URL is required")
	}
	u, err := url.Parse(baseURL)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return fmt.Errorf("invalid base URL: %s", baseURL)
	}
	return nil
}

// newClient loads a cached token and creates an Axosoft API client.
// ClientID is not needed here — it was only required during login.
func newClient(options *Options) (*axosoft.Client, error) {
	if err := validateBaseURL(options.BaseURL); err != nil {
		return nil, err
	}

	token, err := auth.Load(options.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("loading token: %w", err)
	}
	if token == nil {
		return nil, fmt.Errorf("no token found for %s — run 'axosoft-cli login' first", options.BaseURL)
	}

	return axosoft.New(options.BaseURL, token.AccessToken), nil
}
