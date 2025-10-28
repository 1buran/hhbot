package views

import (
	"strings"

	"encoding/json"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type sourceModel struct {
	content      []string
	scroll, w, h int
}

func (m sourceModel) Init() tea.Cmd { return nil }

func (m sourceModel) View() string {
	lines := len(m.content)
	if lines <= m.h { // all content fits the one screen without scroll
		return strings.Join(m.content, "\n")
	}
	return strings.Join(m.content[m.scroll:m.scroll+m.h-1], "\n")
}

func (m sourceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, events.Quit
		case "up", "k":
			if m.scroll > 0 {
				m.scroll--
			}
		case "down", "j":
			if len(m.content)-m.scroll > m.h {
				m.scroll++
			}
		}
	}
	return m, nil
}

func NewSource(v dto.Vacancy, w, h int) (*sourceModel, error) {
	b, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return nil, err
	}
	content := fmt.Sprintf(
		"**Исходное содержание объекта вакансия:**\n\n ```json\n\n%s\n\n```", b)

	renderedContent, err := glamourRenderer.Render(content)
	if err != nil {
		return nil, fmt.Errorf("glamour.Render failure: %w", err)
	}

	return &sourceModel{content: strings.Split(renderedContent, "\n"), w: w, h: h}, nil
}
