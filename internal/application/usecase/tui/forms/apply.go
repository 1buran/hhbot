package forms

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
)

type applyModel struct {
	action           func(message string) tea.Cmd
	textarea         textarea.Model
	greeting         string
	cursor, elements int
}

func (m applyModel) Init() tea.Cmd { return m.textarea.Focus() }

func (m applyModel) View() string {
	var okButton, cancelButton string

	switch m.cursor {
	case 0: // text input in focus, buttons inactive
		okButton = buttonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("No")
	case 1:
		okButton = activeButtonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("No")
	case 2:
		okButton = buttonStyle.Render("Yes")
		cancelButton = activeButtonStyle.Render("No")
	}

	question := lipgloss.NewStyle().Width(formInnerWidth).Align(lipgloss.Center).Render(m.greeting)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	textarea := lipgloss.NewStyle().Margin(1).Render(m.textarea.View())
	ui := lipgloss.JoinVertical(lipgloss.Center, question, textarea, buttons)

	return lipgloss.Place(formWidth, 9,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(subtle),
	)

}

func (m applyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n", "q", "esc":
			return m, events.Quit
		case "y":
			return m, m.action(m.textarea.Value())
		case "tab":
			m.cursor++
			m.cursor = m.cursor % m.elements
			if m.cursor == 0 {
				m.textarea.Focus()
			} else {
				m.textarea.Blur()
			}
		}
	}

	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func NewApplyForm(vacancyID, resumeID, text string, client *apiclient.ApiClient) applyModel {
	action := func(message string) tea.Cmd {
		return func() tea.Msg {
			res, err := client.Apply(vacancyID, resumeID, message)
			if err != nil {
				return events.NewMessage(events.ErrorMessage, err.Error())
			}
			switch res.Code {
			case 201:
				return events.Quit()
			case 303:
				return events.NewMessage(events.Information,
					fmt.Sprintf("Перейдите на: %s", res.Headers.Get("Location")))
			case 400, 403:
				return events.NewMessage(events.ErrorMessage, res.Error)
			}
			return events.NewMessage(events.ErrorMessage, fmt.Sprintf("%+v", res))
		}
	}

	ta := textarea.New()
	ta.Placeholder = "Интересно!"
	ta.SetWidth(formWidth)
	ta.SetHeight(5)
	ta.FocusedStyle.Text = textareaStyle
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false

	return applyModel{
		action:   action,
		greeting: text,
		elements: 3,
		textarea: ta,
	}
}
