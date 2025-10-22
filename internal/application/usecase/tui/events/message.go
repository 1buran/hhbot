package events

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	Information = iota // message type(level)
	ErrorMessage
)

var (
	informationStyle = lipgloss.NewStyle().Width(60).Foreground(lipgloss.Color("47"))
	errorStyle       = lipgloss.NewStyle().Width(60).Foreground(lipgloss.Color("#FF6863"))
)

type Message struct {
	lvl         int
	txt, status string
}

func (m Message) Level() int     { return m.lvl }
func (m Message) Status() string { return m.status }

func (m Message) Init() tea.Cmd { return nil }

func (m Message) View() string {
	switch m.lvl {
	case Information:
		return informationStyle.Render(m.txt)
	case ErrorMessage:
		return errorStyle.Render(m.txt)
	}
	return m.txt
}

func (m Message) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, Quit
		}
	}
	return m, nil
}

func NewMessage(level int, text, status string) Message {
	return Message{lvl: level, txt: text, status: status}
}
