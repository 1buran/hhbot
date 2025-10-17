package events

import (
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type Source struct{ Vacancy dto.Vacancy }

func NewSource(v dto.Vacancy) tea.Cmd {
	return func() tea.Msg {
		return Source{v}
	}
}
