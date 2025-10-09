package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	AuthURL  = "https://hh.ru/oauth/authorize"
	TokenURL = "https://hh.ru/oauth/token"

	redirectURI           = "http://localhost:8080/callback"
	callbackServerAddress = "localhost:8080"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type OAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string

	callbackServer *http.Server
	authCodeChan   chan string
	authStateChan  chan string
}

// GetAuthURL returns the authorization URL with state parameter
func (cfg OAuthClient) GetAuthURL(state string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", cfg.ClientID)
	params.Set("state", state)
	params.Set("redirect_uri", cfg.RedirectURI)

	return fmt.Sprintf("%s?%s", AuthURL, params.Encode())
}

// ExchangeCode exchanges authorization code for access token
func (cfg *OAuthClient) ExchangeCode(code, state string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", cfg.RedirectURI)
	data.Set("state", state)

	req, err := http.NewRequest("POST", TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "HHBot/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed: %d - %s", resp.StatusCode, body)
	}

	var token TokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}

	return &token, nil
}

func (cfg *OAuthClient) callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if code == "" {
		http.Error(w, "No authorization code received", http.StatusBadRequest)
		return
	}

	if state == "" {
		http.Error(w, "No state parameter received", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, `<html>
<head><title>Authorization Successful</title></head>
<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px; background-color: #060707;">
	<h1 style="color: #4CAF50;">✓ Authorization successful!</h1>
	<p style="color: #c9d6d9;">You can close this window and return to the terminal.</p>
</body>
</html>`)

	cfg.authCodeChan <- code
	cfg.authStateChan <- state
}

// Start callback server: it listen on port 8080 and handle only one route /callback.
func (cfg *OAuthClient) StartListeningCallback() {
	http.HandleFunc("/callback", cfg.callbackHandler)

	cfg.callbackServer = &http.Server{Addr: callbackServerAddress}
	go func() {
		fmt.Println("Starting local server on", cfg.callbackServer.Addr)
		if err := cfg.callbackServer.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("Server error:", err)
		}
	}()
}

// Stop callback server.
func (cfg *OAuthClient) StopListeningCallback() {
	if err := cfg.callbackServer.Close(); err != nil {
		fmt.Println(err)
	}
}

// Authentification progress.
func (cfg *OAuthClient) Authenticate() (string, error) {
	// Generate state for CSRF protection
	state, err := GenerateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	expectedState := state

	cfg.StartListeningCallback()
	defer cfg.StopListeningCallback()

	// Get authorization URL with state
	authURL := cfg.GetAuthURL(state)

	fmt.Printf("\nPlease visit this URL to authorize the application:\n\n%s\n\n", authURL)
	fmt.Println("Waiting for authorization...")

	// Wait for authorization code and state
	code := <-cfg.authCodeChan
	receivedState := <-cfg.authStateChan

	// Verify state to prevent CSRF attacks
	if receivedState != expectedState {
		return "", fmt.Errorf("state mismatch: possible CSRF attack")
	}

	// Exchange code for token
	token, err := cfg.ExchangeCode(code, state)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func NewOAuthClient(clientId, clientSecret string) *OAuthClient {
	return &OAuthClient{
		ClientID:      clientId,
		ClientSecret:  clientSecret,
		RedirectURI:   redirectURI,
		authCodeChan:  make(chan string),
		authStateChan: make(chan string),
	}
}
