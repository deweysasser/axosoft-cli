package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- toToken ---

func TestToToken_WithExpiry(t *testing.T) {
	tr := &tokenResponse{AccessToken: "abc", ExpiresIn: 3600}
	before := time.Now()
	token := tr.toToken("https://example.axosoft.com")
	after := time.Now()

	assert.Equal(t, "abc", token.AccessToken)
	assert.Equal(t, "https://example.axosoft.com", token.BaseURL)
	assert.True(t, token.ExpiresAt.After(before.Add(3599*time.Second)))
	assert.True(t, token.ExpiresAt.Before(after.Add(3601*time.Second)))
}

func TestToToken_ZeroExpiresIn(t *testing.T) {
	tr := &tokenResponse{AccessToken: "abc", ExpiresIn: 0}
	token := tr.toToken("https://example.axosoft.com")
	assert.Equal(t, "abc", token.AccessToken)
	assert.Equal(t, "https://example.axosoft.com", token.BaseURL)
	assert.True(t, token.ExpiresAt.IsZero(), "zero ExpiresIn means token never expires")
}

// --- readLimitedErrBody ---

func TestReadLimitedErrBody_ShortBody(t *testing.T) {
	result := readLimitedErrBody(strings.NewReader("error message"))
	assert.Equal(t, "error message", result)
}

func TestReadLimitedErrBody_TruncatesAtLimit(t *testing.T) {
	large := strings.Repeat("x", 5000)
	result := readLimitedErrBody(strings.NewReader(large))
	assert.Len(t, result, maxErrBodyBytes)
}

// --- generateState ---

func TestGenerateState_UniqueEachCall(t *testing.T) {
	s1, err1 := generateState()
	s2, err2 := generateState()
	require.NoError(t, err1)
	require.NoError(t, err2)
	assert.NotEqual(t, s1, s2)
	assert.Len(t, s1, 32) // 16 bytes hex-encoded
	assert.Len(t, s2, 32)
}

// --- exchangeCode ---

func TestExchangeCode_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/oauth2/token", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "authorization_code", r.Form.Get("grant_type"))
		assert.Equal(t, "test-code", r.Form.Get("code"))

		json.NewEncoder(w).Encode(tokenResponse{
			AccessToken: "new-access-token",
			ExpiresIn:   3600,
		})
	}))
	defer srv.Close()

	token, err := exchangeCode(srv.URL, "client-id", "test-code", "http://localhost/callback")
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", token.AccessToken)
	assert.False(t, token.ExpiresAt.IsZero())
}

func TestExchangeCode_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid_grant"))
	}))
	defer srv.Close()

	_, err := exchangeCode(srv.URL, "client-id", "bad-code", "http://localhost/callback")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Contains(t, err.Error(), "invalid_grant")
}

func TestExchangeCode_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not-json{{{"))
	}))
	defer srv.Close()

	_, err := exchangeCode(srv.URL, "client-id", "code", "http://localhost/callback")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing token response")
}

// --- PasswordAuth ---

func TestPasswordAuth_Success(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, r.ParseForm())
		assert.Equal(t, "password", r.Form.Get("grant_type"))
		assert.Equal(t, "testuser", r.Form.Get("username"))

		json.NewEncoder(w).Encode(tokenResponse{
			AccessToken: "pw-token",
			ExpiresIn:   7200,
		})
	}))
	defer srv.Close()

	err := PasswordAuth(srv.URL, "client-id", "secret", "testuser", "pass")
	require.NoError(t, err)

	token, err := Load(srv.URL)
	require.NoError(t, err)
	require.NotNil(t, token)
	assert.Equal(t, "pw-token", token.AccessToken)
}

func TestPasswordAuth_Unauthorized(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
	}))
	defer srv.Close()

	err := PasswordAuth(srv.URL, "client-id", "secret", "user", "wrongpass")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "401")
}
