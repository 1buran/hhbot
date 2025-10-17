package views

import (
	tea "github.com/charmbracelet/bubbletea"

	"fmt"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type listvacanciesModel struct {
	action func(vacancyID string) tea.Cmd

	vacancies            []dto.Vacancy
	cursor, scroll, w, h int

	hhdict dto.Dictionary
}

func (m listvacanciesModel) Init() tea.Cmd { return nil }

func (m listvacanciesModel) View() string {
	return renderVacancy(m.cursor, m.vacancies[m.cursor], m.hhdict)
}

func (m listvacanciesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.action(m.vacancies[m.cursor].Id)
		case "a":
			v := m.vacancies[m.cursor]
			return m, events.NewApply(m.cursor, v.Id, v.Name, v.Salary_range.String(),
				v.Employer.Name, v.Archived, v.Response_letter_required)
		case "ctrl+c", "q":
			return m, tea.Quit
		case "s":
			return m, events.NewSource(m.vacancies[m.cursor])
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.vacancies)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func NewListVacancies(
	vacancies []dto.Vacancy, dict dto.Dictionary, client *apiclient.ApiClient,
) *listvacanciesModel {
	action := func(vacancyID string) tea.Cmd {
		return func() tea.Msg {
			vac, err := client.GetVacancy(vacancyID)
			if err != nil {
				return events.NewMessage(events.ErrorMessage,
					fmt.Errorf("Get vacancy failure: %w", err).Error())
			}
			return events.NewShowVacancy(vac.Name, vac.Salary_range.String(),
				vac.Key_skills.String(), vac.Description, vac.Employer.Name,
				vac.Work_format.String())
		}
	}

	return &listvacanciesModel{vacancies: vacancies, hhdict: dict, action: action}
}
