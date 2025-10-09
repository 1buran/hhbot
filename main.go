package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/auth"
)

var (
	vacancyCard   = lipgloss.NewStyle().PaddingLeft(4)
	redflagStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6863"))
	appliedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#80EF80"))
	invitedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	questionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00F0FF"))
	favoriteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7518"))

	blacklistedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFEF00"))
	greenlightStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("121"))

	salaryStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#dddddd")).Bold(true)
	companyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Italic(true)
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

	relations, ok := dict.Get("vacancy_relation")
	if !ok {
		fmt.Printf("Warning: vacancy_relation not found! dict: %+v\n", dict)
	}
	fmt.Printf("Виды связи: %+v\n", relations)

	fmt.Println("\n=== Searching for vacancies ===")
	vacancies, err := client.SearchVacancies("Golang developer", 0)
	if err != nil {
		fmt.Println("client.SearchVacancies failure:", err)
		return
	}

	for i, v := range vacancies {
		if v.Archived {
			continue
		}
		var rel []string
		for _, k := range v.Relations {
			v := relations.Get(k)
			if v != "" {
				switch k {
				case "favorite":
					v = favoriteStyle.Render(v)
				case "got_response":
					v = appliedStyle.Render(v)
				case "got_rejection":
					v = redflagStyle.Render(v)
				case "got_invitation":
					v = invitedStyle.Render(v)
				case "blacklisted":
					v = blacklistedStyle.Render(v)
				case "got_question":
					v = questionStyle.Render(v)
				}
				rel = append(rel, v)
			}
		}

		var experience string
		switch v.Experience.Id {
		case "moreThan6":
			experience = redflagStyle.Render(v.Experience.Name)
		case "between1And3", "noExperience":
			experience = greenlightStyle.Render(v.Experience.Name)
		default:
			experience = v.Experience.Name
		}

		fmt.Printf("%d. %s / %s\n", i+1, titleStyle.Render(v.Name),
			salaryStyle.Render(v.Salary_range.String()))
		fmt.Println(
			vacancyCard.Render(
				fmt.Sprintf(
					"Компания: %s\nОпыт: %s\nСвязь: %s\nФормат работы: %s\n"+
						"Откликнулось: %d / %d\nОпубликовано: %s\nСоздано: %s\n%s\n\n",
					companyStyle.Render(v.Employer.Name),
					experience,
					strings.Join(rel, ", "),
					v.Work_format.String(),
					v.Counters.Responses, v.Counters.Total_responses,
					v.Published_at.Format(time.DateOnly),
					v.Created_at.Format(time.DateOnly),
					v.Alternate_url,
				),
			),
		)
	}
}
