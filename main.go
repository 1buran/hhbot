package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/auth"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type model struct {
	vacancies []dto.Vacancy // load vacancies from hh.ru
	cursor    int           // active vacancy
	hhdict    dto.Dictionary
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) View() string {
	return renderVacancy(m.cursor, m.vacancies[m.cursor], m.hhdict)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.vacancies)-1 {
				m.cursor++
			}

		// Enter to vacancy
		case "enter", " ":
			// todo load vacancy and show it in alternate view
		}

	}
	return m, nil
}

func initialModel(vacancies []dto.Vacancy, dict dto.Dictionary) model {
	return model{vacancies: vacancies, hhdict: dict}
}

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

	p := tea.NewProgram(initialModel(vacancies, dict))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
