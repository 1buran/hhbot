package styles

import "github.com/charmbracelet/lipgloss"

var (
	VacancyCard = lipgloss.NewStyle().PaddingLeft(4)

	Redflag = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6863"))
	Applied = lipgloss.NewStyle().Foreground(lipgloss.Color("#80EF80"))
	Invited = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))

	Question = lipgloss.NewStyle().Foreground(lipgloss.Color("#00F0FF"))
	Favorite = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7518"))

	Blacklisted = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFEF00"))
	Greenlight  = lipgloss.NewStyle().Foreground(lipgloss.Color("121"))

	Title   = lipgloss.NewStyle().Foreground(lipgloss.Color("#dddddd")).Bold(true)
	Salary  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	Company = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Italic(true)

	Information = lipgloss.NewStyle().Width(60).Foreground(lipgloss.Color("47"))

	Action       = lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
	ActionPrompt = lipgloss.NewStyle().Inherit(Action).SetString("⟩")
	ActionInput  = lipgloss.NewStyle().Blink(true).Inherit(Action)
)
