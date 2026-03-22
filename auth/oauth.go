package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
)

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Authenticate runs the OAuth browser flow and returns a saved token.
func Authenticate(baseURL, clientID string) error {
	baseURL = strings.TrimRight(baseURL, "/")

	// Start local HTTP server on a random port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return fmt.Errorf("starting callback server: %w", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port
	redirectURI := fmt.Sprintf("http://localhost:%d/callback", port)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error")
			if errMsg == "" {
				errMsg = "no code in callback"
			}
			errCh <- fmt.Errorf("OAuth error: %s", errMsg)
			fmt.Fprintf(w, "<html><body><h1>Authentication failed</h1><p>%s</p></body></html>", errMsg)
			return
		}
		codeCh <- code
		fmt.Fprint(w, "<html><body><h1>Authentication successful!</h1><p>You can close this window.</p></body></html>")
	})

	srv := &http.Server{Handler: mux}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	defer srv.Shutdown(context.Background())

	// Build authorization URL and open browser
	authURL := fmt.Sprintf("%s/auth?response_type=code&client_id=%s&redirect_uri=%s&scope=read%%20write",
		baseURL,
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
	)

	log.Info().Str("url", authURL).Msg("Opening browser for authentication")
	if err := browser.OpenURL(authURL); err != nil {
		log.Warn().Err(err).Msg("Could not open browser automatically")
		fmt.Printf("Please open this URL in your browser:\n%s\n", authURL)
	}

	// Wait for callback
	var code string
	select {
	case code = <-codeCh:
	case err := <-errCh:
		return err
	case <-time.After(5 * time.Minute):
		return fmt.Errorf("authentication timed out")
	}

	// Exchange code for token
	token, err := exchangeCode(baseURL, clientID, code, redirectURI)
	if err != nil {
		return err
	}

	return Save(token)
}

// PasswordAuth authenticates using the username/password grant type.
func PasswordAuth(baseURL, clientID, clientSecret, username, password string) error {
	baseURL = strings.TrimRight(baseURL, "/")

	params := url.Values{
		"grant_type":    {"password"},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {username},
		"password":      {password},
		"scope":         {"read write"},
	}

	tokenURL := baseURL + "/api/oauth2/token?" + params.Encode()

	resp, err := http.Get(tokenURL)
	if err != nil {
		return fmt.Errorf("requesting token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("parsing token response: %w", err)
	}

	token := &Token{
		AccessToken: tokenResp.AccessToken,
		BaseURL:     baseURL,
	}
	if tokenResp.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	log.Info().Msg("Password authentication successful")
	return Save(token)
}

func exchangeCode(baseURL, clientID, code, redirectURI string) (*Token, error) {
	tokenURL := baseURL + "/api/oauth2/token"

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"client_id":    {clientID},
		"redirect_uri": {redirectURI},
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d", resp.StatusCode)
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}

	token := &Token{
		AccessToken: tokenResp.AccessToken,
		BaseURL:     baseURL,
	}
	if tokenResp.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	log.Info().Msg("Authentication successful")
	return token, nil
}
