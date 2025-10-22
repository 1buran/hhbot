package events

import tea "github.com/charmbracelet/bubbletea"

type StatusBarNotifyLevel int

const (
	StatusBarNotifyLevelDefault StatusBarNotifyLevel = iota
	StatusBarNotifyLevelWarning
	StatusBarNotifyLevelAlert
)

type StatusBarNotify struct {
	Text  string
	Level StatusBarNotifyLevel
}

func NewStatusBarNotifyMsg(txt string, lvl StatusBarNotifyLevel) StatusBarNotify {
	return StatusBarNotify{Text: txt, Level: lvl}
}

func NewStatusBarNotifyCmd(txt string, lvl StatusBarNotifyLevel) tea.Cmd {
	return func() tea.Msg {
		return StatusBarNotify{Text: txt, Level: lvl}
	}
}
