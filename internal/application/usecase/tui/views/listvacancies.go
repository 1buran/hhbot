package views

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/actions"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/state"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type listvacanciesModel struct {
	action func(vacancyID string) tea.Cmd

	vacancies []dto.Vacancy
	cursor    int

	hhdict dto.Dictionary
	state  *state.Facade
}

func (m listvacanciesModel) Init() tea.Cmd { return nil }

func (m listvacanciesModel) View() string {
	if len(m.vacancies) > 0 {
		v := m.vacancies[m.cursor]
		_, vb := m.state.IsVacancyBlacklisted(v.Id)
		_, cb := m.state.IsCompanyBlacklisted(v.Employer.Id)
		return renderVacancy(m.cursor, v, m.hhdict, vb || cb)
	}
	return styles.Information.Render("Ничего не найдено! Повторите поиск!")
}

func (m listvacanciesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, m.action(m.vacancies[m.cursor].Id)
		case "ctrl+f":
			return m, events.NewUserInput()
		case "b":
			return m, actions.NewBlacklisted(m.vacancies[m.cursor].Id, "", "поблэклистили")
		case "B":
			return m, actions.NewBlacklisted("", m.vacancies[m.cursor].Employer.Id, "поблэклистили")
		case "u":
			return m, actions.NewUnblacklisted(m.vacancies[m.cursor].Id, "", "разблокировали")
		case "U":
			return m, actions.NewUnblacklisted("", m.vacancies[m.cursor].Employer.Id, "разблокировали")
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
	vacancies []dto.Vacancy, dict dto.Dictionary,
	client *apiclient.ApiClient, state *state.Facade,
) *listvacanciesModel {
	action := func(vacancyID string) tea.Cmd {
		return func() tea.Msg {
			vac, err := client.GetVacancy(vacancyID)
			if err != nil {
				return events.NewMessage(events.ErrorMessage,
					fmt.Errorf("get vacancy failure: %w", err).Error(),
					"Ошибка просмотра вакансии")
			}
			return events.NewShowVacancy(vac.Name, vac.Salary_range.String(),
				vac.Key_skills.String(), vac.Description, vac.Employer.Name,
				vac.Work_format.String())
		}
	}

	return &listvacanciesModel{vacancies: vacancies, hhdict: dict, action: action, state: state}
}
