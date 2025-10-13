package main

import (
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	formInnerWidth, formWidth = 90, 100
)

var (
	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			Margin(1)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94"))

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	textareaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("189"))

	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
)

type Form interface {
	Tab() tea.Cmd
	View() string
	Focused() bool
	Update(msg tea.Msg) tea.Cmd
	Message() string
}

type dialogyesno struct {
	text             string
	textarea         textarea.Model
	cursor, elements int
	focused          bool // text area is active, all type events must be forward to it
}

func (d *dialogyesno) Tab() tea.Cmd {
	d.cursor++
	d.cursor = d.cursor % d.elements
	if d.cursor == 0 {
		d.focused = true
		return d.textarea.Focus()
	}
	d.focused = false
	d.textarea.Blur()
	return nil
}

func (d dialogyesno) Message() string { return d.textarea.Value() }
func (d dialogyesno) Focused() bool   { return d.focused }
func (d *dialogyesno) Update(msg tea.Msg) (cmd tea.Cmd) {
	d.textarea, cmd = d.textarea.Update(msg)
	return cmd
}

func (d dialogyesno) View() string {
	var okButton, cancelButton string

	switch d.cursor {
	case 0:
		okButton = buttonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("No")
	case 1:
		okButton = activeButtonStyle.Render("Yes")
		cancelButton = buttonStyle.Render("No")
	case 2:
		okButton = buttonStyle.Render("Yes")
		cancelButton = activeButtonStyle.Render("No")
	}

	question := lipgloss.NewStyle().Width(formInnerWidth).Align(lipgloss.Center).Render(d.text)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	textarea := lipgloss.NewStyle().Margin(1).Render(d.textarea.View())
	ui := lipgloss.JoinVertical(lipgloss.Center, question, textarea, buttons)

	return lipgloss.Place(formWidth, 9,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}

func NewDialogYesNo(text string) Form {
	ta := textarea.New()
	ta.Placeholder = "Интересно!"
	ta.SetWidth(formWidth)
	ta.SetHeight(5)
	ta.FocusedStyle.Text = textareaStyle
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.ShowLineNumbers = false
	return &dialogyesno{text: text, elements: 3, textarea: ta} // todo: add text input (message)
}
