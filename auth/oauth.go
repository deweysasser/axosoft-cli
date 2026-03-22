package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/browser"
	"github.com/rs/zerolog/log"
)

const (
	authHTTPTimeout = 15 * time.Second
	maxErrBodyBytes = 4096
)

// authClient is a shared HTTP client with timeouts for OAuth token requests.
var authClient = &http.Client{Timeout: authHTTPTimeout}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func (tr *tokenResponse) toToken(baseURL string) *Token {
	token := &Token{
		AccessToken: tr.AccessToken,
		BaseURL:     baseURL,
	}
	if tr.ExpiresIn > 0 {
		token.ExpiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	}
	return token
}

func readLimitedErrBody(body io.Reader) string {
	data, _ := io.ReadAll(io.LimitReader(body, maxErrBodyBytes))
	return string(data)
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating state: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// Authenticate runs the OAuth browser flow and returns a saved token.
func Authenticate(baseURL, clientID string) error {
	baseURL = strings.TrimRight(baseURL, "/")

	state, err := generateState()
	if err != nil {
		return err
	}

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
		if r.URL.Query().Get("state") != state {
			http.Error(w, "invalid state", http.StatusBadRequest)
			errCh <- fmt.Errorf("OAuth CSRF: state mismatch")
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			errMsg := r.URL.Query().Get("error")
			if errMsg == "" {
				errMsg = "no code in callback"
			}
			errCh <- fmt.Errorf("OAuth error: %s", errMsg)
			fmt.Fprintf(w, "<html><body><h1>Authentication failed</h1><p>%s</p></body></html>", html.EscapeString(errMsg))
			return
		}
		codeCh <- code
		fmt.Fprint(w, "<html><body><h1>Authentication successful!</h1><p>You can close this window.</p></body></html>")
	})

	srv := &http.Server{
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Warn().Err(err).Msg("Error shutting down callback server")
		}
	}()

	// Build authorization URL and open browser
	authURL := fmt.Sprintf("%s/auth?response_type=code&client_id=%s&redirect_uri=%s&scope=read%%20write&state=%s",
		baseURL,
		url.QueryEscape(clientID),
		url.QueryEscape(redirectURI),
		url.QueryEscape(state),
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

	tokenURL := baseURL + "/api/oauth2/token"

	resp, err := authClient.PostForm(tokenURL, params)
	if err != nil {
		return fmt.Errorf("requesting token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed (status %d): %s", resp.StatusCode, readLimitedErrBody(resp.Body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("parsing token response: %w", err)
	}

	log.Info().Msg("Password authentication successful")
	return Save(tokenResp.toToken(baseURL))
}

func exchangeCode(baseURL, clientID, code, redirectURI string) (*Token, error) {
	tokenURL := baseURL + "/api/oauth2/token"

	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"client_id":    {clientID},
		"redirect_uri": {redirectURI},
	}

	resp, err := authClient.PostForm(tokenURL, data)
	if err != nil {
		return nil, fmt.Errorf("exchanging code for token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed (status %d): %s", resp.StatusCode, readLimitedErrBody(resp.Body))
	}

	var tokenResp tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("parsing token response: %w", err)
	}

	log.Info().Msg("Authentication successful")
	return tokenResp.toToken(baseURL), nil
}
