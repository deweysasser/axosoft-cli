package program

import (
	"fmt"

	"github.com/deweysasser/axosoft-cli/auth"
)

// LoginCmd handles authentication with Axosoft
type LoginCmd struct {
	Username     string `help:"Axosoft username (for password grant flow)" env:"AXOSOFT_USERNAME"`
	Password     string `help:"Axosoft password — prefer AXOSOFT_PASSWORD env var to avoid process listing exposure" env:"AXOSOFT_PASSWORD"`
	ClientSecret string `help:"Axosoft client secret — prefer AXOSOFT_CLIENT_SECRET env var" env:"AXOSOFT_CLIENT_SECRET"`
}

func (cmd *LoginCmd) Run(options *Options) error {
	if err := validateBaseURL(options.BaseURL); err != nil {
		return err
	}

	// Use password grant if username is provided
	if cmd.Username != "" {
		if cmd.Password == "" {
			return fmt.Errorf("--password is required when using --username")
		}
		if options.ClientID == "" {
			return fmt.Errorf("--client-id or AXOSOFT_CLIENT_ID is required for password auth (create an OAuth client in Axosoft > Tools > System Settings > API Keys)")
		}
		if cmd.ClientSecret == "" {
			return fmt.Errorf("--client-secret is required for password auth")
		}
		return auth.PasswordAuth(options.BaseURL, options.ClientID, cmd.ClientSecret, cmd.Username, cmd.Password)
	}

	if options.ClientID == "" {
		return fmt.Errorf("--client-id or AXOSOFT_CLIENT_ID is required")
	}
	return auth.Authenticate(options.BaseURL, options.ClientID)
}
