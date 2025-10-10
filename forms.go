package main

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2)
		// Underline(true)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
)

func DialogYesNo(text string) string {
	okButton := activeButtonStyle.Render("Yes")
	cancelButton := buttonStyle.Render("No")

	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(text)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, question, buttons)

	width := 80

	return lipgloss.Place(width, 9,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(subtle),
	)
}
