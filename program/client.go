package program

import (
	"fmt"
	"net/url"

	"github.com/deweysasser/axosoft-cli/auth"
	"github.com/deweysasser/axosoft-cli/axosoft"
)

// newClient loads a cached token and creates an Axosoft API client.
// ClientID is not needed here — it was only required during login.
func newClient(options *Options) (*axosoft.Client, error) {
	if options.BaseURL == "" {
		return nil, fmt.Errorf("--base-url or AXOSOFT_URL is required")
	}

	u, err := url.Parse(options.BaseURL)
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") {
		return nil, fmt.Errorf("invalid base URL: %s", options.BaseURL)
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
