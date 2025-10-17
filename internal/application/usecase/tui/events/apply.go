package events

import tea "github.com/charmbracelet/bubbletea"

type Apply struct {
	VacancyNumber       int
	Title, Salary       string
	VacancyID, Employer string

	Archived, ResponseRequired bool
}

func NewApply(vacNumber int, vacID, vacTitle, vacSalary, employer string,
	archived, respRequired bool) tea.Cmd {
	return func() tea.Msg {
		return Apply{
			VacancyNumber:    vacNumber,
			VacancyID:        vacID,
			Title:            vacTitle,
			Salary:           vacSalary,
			Archived:         archived,
			Employer:         employer,
			ResponseRequired: respRequired,
		}
	}
}
