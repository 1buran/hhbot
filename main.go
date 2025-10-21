package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/forms"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/state"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/auth"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
	"gitlab.com/1buran/hhbot/internal/infrastructure/cache"
)

func main() {
	clientID := os.Getenv("HH_CLIENT_ID")
	clientSecret := os.Getenv("HH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set HH_CLIENT_ID and HH_CLIENT_SECRET in .env file")
	}

	resumeID, ok := os.LookupEnv("RESUME_ID")
	if !ok {
		fmt.Println(styles.Redflag.Render("Not found env var: RESUME_ID"))
		return
	}

	var (
		err           error
		accessToken   string
		tokenResponse auth.TokenResponse
	)

	if err := cache.LoadToken(&tokenResponse); err != nil {
		fmt.Println("access token load from cache failure:", err.Error())
	} else {
		accessToken = tokenResponse.AccessToken
	}

	if accessToken == "" {
		oauthClient := auth.NewOAuthClient(clientID, clientSecret)

		// Start OAuth flow
		fmt.Println("Starting OAuth authentication...")

		tokenResponse, err := oauthClient.Authenticate()
		if err != nil {
			log.Fatal("Authentication failed:", err)
		}
		accessToken = tokenResponse.AccessToken

		fmt.Println(styles.Information.Render("\n✓ Authentication successful!"))
		if err := cache.SaveToken(tokenResponse); err != nil {
			fmt.Println(styles.Redflag.Render(err.Error()))
		}
	}

	// Create API client
	client := apiclient.NewApiClient(accessToken)

	// Example: Search for vacancies
	var dict dto.Dictionary
	fmt.Println(styles.ActionPrompt.String(), styles.ActionInput.Render("Get dictionary..."))
	if dict, err = client.GetDictionary(); err != nil {
		fmt.Println(styles.Redflag.Render("client.GetDictionary failure:", err.Error()))
		return
	}
	fmt.Println(styles.Information.Render("✓ hh.ru dictionary loaded successful"))
	time.Sleep(500 * time.Millisecond)

	sts := state.NewState()
	var prevState state.State
	if err := cache.LoadState(&prevState); err != nil {
		// no previous state available, it's ok
	} else {
		sts.Init(&prevState)
	}

	ticker := time.NewTicker(time.Second)
	done := make(chan struct{})
	defer close(done)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if sts.IsDirty() {
					if err := cache.SaveState(sts.ExportState()); err != nil {
						tea.Println(styles.Redflag.Render(err.Error()))
					} else {
						sts.SetClean()
					}
				}
			}
		}
	}()

	p := tea.NewProgram(
		tui.InitialModel(forms.NewInput(client), client, sts, dict, resumeID),
		tea.WithAltScreen(),
		//		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Run program failure: %v", err)
		done <- struct{}{}
		os.Exit(1)
	}
	done <- struct{}{}
}
