package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

	relations, ok := dict.Get("vacancy_relation")
	if !ok {
		fmt.Printf("Warning: vacancy_relation not found! dict: %+v\n", dict)
	}
	fmt.Printf("Виды связи: %+v\n", relations)

	fmt.Println("\n=== Searching for vacancies ===")
	vacancies, err := client.SearchVacancies("Golang developer", 350_000)
	if err != nil {
		fmt.Println("client.SearchVacancies failure:", err)
		return
	}

	for i, v := range vacancies {
		fmt.Printf("%d. %s\n", i+1, v.Name)
		fmt.Printf("    %s\n", v.Salary_range)
		fmt.Printf("    Компания: %s\n", v.Employer.Name)
		fmt.Printf("    Опыт: %s\n", v.Experience.Name)

		var rel []string
		for _, k := range v.Relations {
			v := relations.Get(k)
			if v != "" {
				rel = append(rel, v)
			}
		}

		fmt.Printf("    Связь: %s\n", strings.Join(rel, ", "))
		fmt.Printf("    Формат работы: %s\n", v.Work_format.String())
		fmt.Printf("    Откликнулось: %d / %d\n", v.Counters.Responses, v.Counters.Total_responses)
		fmt.Printf("    Опубликовано: %s\n", v.Published_at.Format(time.DateOnly))
		fmt.Printf("    Создано: %s\n", v.Created_at.Format(time.DateOnly))
		fmt.Printf("    %s\n\n", v.Alternate_url)
	}
}
