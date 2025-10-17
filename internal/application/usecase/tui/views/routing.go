package views

import tea "github.com/charmbracelet/bubbletea"

type History struct {
	views []tea.Model
}

func (d *History) Save(v tea.Model) {
	if v != nil {
		d.views = append(d.views, v)
	}
}

func (d *History) Back() (v tea.Model) {
	if len(d.views) > 0 {
		v, d.views = d.views[len(d.views)-1], d.views[:len(d.views)-1]
	}
	return
}

func NewHistory() *History { return &History{} }
