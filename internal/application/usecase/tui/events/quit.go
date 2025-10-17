package events

import tea "github.com/charmbracelet/bubbletea"

type QuitFromInnerView struct{}

func Quit() tea.Msg { return QuitFromInnerView{} }
