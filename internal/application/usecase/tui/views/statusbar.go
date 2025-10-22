package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
)

type statusbarModel struct {
	txt string
	w   int
	lvl events.StatusBarNotifyLevel
}

var (
	statusbarStyleDefault = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(1).
				Align(lipgloss.Left).
				Background(lipgloss.Color("235")).
				Foreground(lipgloss.Color("189"))

	statusbarStyleWarning = lipgloss.NewStyle().
				MarginLeft(1).
				MarginRight(1).
				Foreground(lipgloss.Color("#ffff00")).
				Inherit(statusbarStyleDefault)

	statusbarStyleAlert = lipgloss.NewStyle().
				Bold(true).
				MarginLeft(1).
				MarginRight(1).
				Foreground(lipgloss.Color("#ff000d")).
				Inherit(statusbarStyleDefault)
)

func (m statusbarModel) Init() tea.Cmd { return nil }

func (m statusbarModel) View() string {
	if m.txt == "" {
		return ""
	}

	switch m.lvl {
	case events.StatusBarNotifyLevelWarning:
		return statusbarStyleWarning.Width(m.w).Render(m.txt)
	case events.StatusBarNotifyLevelAlert:
		return statusbarStyleAlert.Width(m.w).Render(m.txt)
	}
	return statusbarStyleDefault.Width(m.w).Render(m.txt)
}

func (m statusbarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w = msg.Width
	case events.StatusBarNotify:
		m.lvl, m.txt = msg.Level, msg.Text
	}

	return m, nil
}

func NewStatusBar() statusbarModel { return statusbarModel{} }

func NewStatusBarWithContext(lvl events.StatusBarNotifyLevel, txt string, width int) statusbarModel {
	return statusbarModel{lvl: lvl, txt: txt, w: width}
}
