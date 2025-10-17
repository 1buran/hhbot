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
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/auth"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
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
		err         error
		accessToken string
	)

	if bytes, err := cacheLoad(".accesstoken"); err != nil {
		fmt.Println("access token load from cache failure:", err.Error())
	} else {
		accessToken = string(bytes)
	}

	if accessToken == "" {
		oauthClient := auth.NewOAuthClient(clientID, clientSecret)

		// Start OAuth flow
		fmt.Println("Starting OAuth authentication...")

		accessToken, err = oauthClient.Authenticate()
		if err != nil {
			log.Fatal("Authentication failed:", err)
		}

		fmt.Println(styles.Information.Render("\n✓ Authentication successful!"))
		if err := cacheSave(".accesstoken", []byte(accessToken)); err != nil {
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

	// fmt.Println("\n=== Searching for vacancies ===")
	// vacancies, err := client.SearchVacancies(
	// 	"NAME:(Golang OR Go NOT (devops OR Яндекс)) AND NOT COMPANY_NAME:(Яндекс OR Yandex)",
	// 	0)
	// if err != nil {
	// 	fmt.Println("client.SearchVacancies failure:", err)
	// 	return
	// }
	// lv := views.NewListVacancies(vacancies, dict, client)

	p := tea.NewProgram(
		tui.InitialModel(forms.NewInput(client), client, dict, resumeID),
		tea.WithAltScreen(),
		//		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
