package actions

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Unblacklisted struct {
	VacancyID, CompanyID, Note string
	Date                       time.Time
}

func NewUnblacklisted(vacancyID, companyID, note string) tea.Cmd {
	return func() tea.Msg {
		return Unblacklisted{
			VacancyID: vacancyID,
			CompanyID: companyID,
			Note:      note,
			Date:      time.Now(),
		}
	}
}
