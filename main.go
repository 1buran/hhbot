package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/1buran/hhbot/internal/infrastructure/apiclient"
)

const (
	redirectPort = "8080"
	redirectURI  = "http://localhost:8080/callback"
)

var (
	authCodeChan  = make(chan string)
	authStateChan = make(chan string)
	oauthConfig   *apiclient.OAuthConfig
	expectedState string
)

func main() {
	// Load credentials from .env file
	loadEnv()

	clientID := os.Getenv("HH_CLIENT_ID")
	clientSecret := os.Getenv("HH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set HH_CLIENT_ID and HH_CLIENT_SECRET in .env file")
	}

	oauthConfig = &apiclient.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}

	// Start OAuth flow
	fmt.Println("=== HH.ru Vacancy Search Bot ===")
	fmt.Println("Starting OAuth authentication...")

	accessToken, err := authenticate()
	if err != nil {
		log.Fatal("Authentication failed:", err)
	}

	fmt.Println("\n✓ Authentication successful!")

	// Create API client
	client := apiclient.NewApiClient(accessToken)

	// Example: Search for vacancies
	fmt.Println("\n=== Searching for vacancies ===")
	vacancies, err := client.SearchVacancies("Golang developer", 350_000)
	if err != nil {
		fmt.Println("client.SearchVacancies failure:", err)
		return
	}
	for i, v := range vacancies {
		fmt.Printf("%d %+v\n", i, v)
	}
}

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // .env file is optional
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}

func authenticate() (string, error) {
	// Generate state for CSRF protection
	state, err := apiclient.GenerateState()
	if err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	expectedState = state

	// Start local server for OAuth callback
	http.HandleFunc("/callback", callbackHandler)

	server := &http.Server{Addr: ":" + redirectPort}

	go func() {
		fmt.Printf("Starting local server on http://localhost:%s\n", redirectPort)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Get authorization URL with state
	authURL := oauthConfig.GetAuthURL(state)

	fmt.Printf("\nPlease visit this URL to authorize the application:\n\n%s\n\n", authURL)
	fmt.Println("Waiting for authorization...")

	// Wait for authorization code and state
	code := <-authCodeChan
	receivedState := <-authStateChan

	// Verify state to prevent CSRF attacks
	if receivedState != expectedState {
		server.Close()
		return "", fmt.Errorf("state mismatch: possible CSRF attack")
	}

	// Shutdown server
	server.Close()

	// Exchange code for token
	token, err := oauthConfig.ExchangeCode(code, state)
	if err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
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
<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px;">
	<h1 style="color: #4CAF50;">✓ Authorization successful!</h1>
	<p>You can close this window and return to the terminal.</p>
</body>
</html>`)

	authCodeChan <- code
	authStateChan <- state
}
