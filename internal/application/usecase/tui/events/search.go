package events

import (
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type NewSearch struct{}

type SearchResults struct {
	Vacancies []dto.Vacancy
}

func NewUserInput() tea.Cmd {
	return func() tea.Msg {
		return NewSearch{}
	}
}

func NewSearchResults(vacancies []dto.Vacancy) SearchResults {
	return SearchResults{vacancies}
}
