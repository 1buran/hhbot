package forms

import (
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
