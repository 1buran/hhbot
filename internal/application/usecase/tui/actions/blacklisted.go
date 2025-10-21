package actions

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Blacklisted struct {
	VacancyID, CompanyID, Note string
	Date                       time.Time
}

func NewBlacklisted(vacancyID, companyID, note string) tea.Cmd {
	return func() tea.Msg {
		return Blacklisted{
			VacancyID: vacancyID,
			CompanyID: companyID,
			Note:      note,
			Date:      time.Now(),
		}
	}
}
