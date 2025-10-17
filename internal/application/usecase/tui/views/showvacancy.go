package views

import (
	"fmt"
	"strings"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/joho/godotenv/autoload"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
)

type showvacancyModel struct {
	content              []string
	cursor, scroll, w, h int
}

func (m showvacancyModel) Init() tea.Cmd { return nil }

func (m showvacancyModel) View() string {
	lines := len(m.content)
	if lines <= m.h { // all content fits the one screen without scroll
		return strings.Join(m.content, "\n")
	}
	return strings.Join(m.content[m.scroll:m.scroll+m.h-1], "\n")
}

func (m showvacancyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func NewShowVacancy(
	title, salary, skills, format, description, employer string,
	w, h int,
) (*showvacancyModel, error) {
	title = fmt.Sprintf("**%s** / `%s` / _%s_\n\n---\n\n", title, salary, employer)
	skills = fmt.Sprintf("**Skills**: %s\n\n", skills)
	format = fmt.Sprintf("**Формат работы**: %s\n\n", format)

	var err error
	description, err = html2md.ConvertString(description)
	if err != nil {
		return nil, err
	}

	renderedContent, err := glamourRenderer.Render(title + skills + format + description)
	if err != nil {
		return nil, err
	}

	return &showvacancyModel{content: strings.Split(renderedContent, "\n"), w: w, h: h}, nil
}
