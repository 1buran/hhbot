package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/auth"
)

func main() {
	clientID := os.Getenv("HH_CLIENT_ID")
	clientSecret := os.Getenv("HH_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		log.Fatal("Please set HH_CLIENT_ID and HH_CLIENT_SECRET in .env file")
	}

	oauthClient := auth.NewOAuthClient(clientID, clientSecret)

	// Start OAuth flow
	fmt.Println("=== HH.ru Vacancy Search Bot ===")
	fmt.Println("Starting OAuth authentication...")

	accessToken, err := oauthClient.Authenticate()
	if err != nil {
		log.Fatal("Authentication failed:", err)
	}

	fmt.Println("\n✓ Authentication successful!")

	// Create API client
	client := apiclient.NewApiClient(accessToken)

	// Example: Search for vacancies
	fmt.Println("Get dictionary...")
	dict, err := client.GetDictionary()
	if err != nil {
		fmt.Println("client.GetDictionary failure:", err)
		return
	}

	fmt.Println("\n=== Searching for vacancies ===")
	vacancies, err := client.SearchVacancies("Golang developer", 0)
	if err != nil {
		fmt.Println("client.SearchVacancies failure:", err)
		return
	}

	// for i, v := range vacancies {
	// 	if v.Archived {
	// 		continue
	// 	}
	// 	fmt.Print(renderVacancy(i, v, dict))
	// }

	p := tea.NewProgram(
		initialModel(client, vacancies, dict),
		tea.WithAltScreen(),
		//		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
