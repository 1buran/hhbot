package forms

import (
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
)

type inputModel struct {
	textinput textinput.Model
	action    func(input string) tea.Cmd
}

func (m inputModel) Init() tea.Cmd { return m.textinput.Focus() }

func (m inputModel) View() string { return m.textinput.View() }

func (m inputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+r":
			m.textinput.Reset()
			m.textinput.Focus()
			return m, nil
		case "enter":
			m.textinput.Blur()
			return m, tea.Batch(m.action(m.textinput.Value()),
				events.NewStatusBarNotifyCmd("Ожидаем ответа от hh.ru...",
					events.StatusBarNotifyLevelWarning))
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
				return events.NewMessage(events.ErrorMessage, err.Error(),
					"Ошибка поиска вакансий")
			}
			if len(vacancies) == 0 {
				return events.NewMessage(events.Information,
					"Ничего не найдено! Повторите поиск!", "Инофрмация")
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
