package forms

import (
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
)

type inputModel struct {
	textinput   textinput.Model
	waitresults bool
	elapsed     int
	action      func(input string) tea.Cmd
}

func (m inputModel) Init() tea.Cmd { return m.textinput.Focus() }

func (m inputModel) View() string {
	s := m.textinput.View()
	if m.waitresults {
		s += "\n\n" + " Waiting response from hh.ru API..."
	}
	return s
}

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			m.textinput.Reset()
			m.textinput.Focus()
			m.waitresults = false
			return m, nil
		case "enter":
			m.textinput.Blur()
			m.waitresults = true
			return m, m.action(m.textinput.Value())
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	m.textinput, cmd = m.textinput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func NewInput(client *apiclient.ApiClient) *inputModel {
	action := func(query string) tea.Cmd {
		return func() tea.Msg {
			vacancies, err := client.SearchVacancies(query, 0)
			if err != nil {
				return events.NewMessage(events.ErrorMessage, err.Error())
			}
			if len(vacancies) == 0 {
				return events.NewMessage(events.Information,
					"Ничего не найдено! Повторите поиск!")
			}
			return events.NewSearchResults(vacancies)
		}
	}

	ti := textinput.New()
	ti.Prompt = "⟩ "
	ti.PromptStyle = styles.Action
	ti.TextStyle = styles.Action
	ti.Placeholder = "Golang"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 20

	return &inputModel{action: action, textinput: ti}
}
