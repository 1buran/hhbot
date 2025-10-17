package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/application/usecase/events"
	"gitlab.com/1buran/hhbot/internal/application/usecase/forms"
	"gitlab.com/1buran/hhbot/internal/application/usecase/views"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

// todo: move styles to separate file
var (
	redflagStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6863"))
	informationStyle = lipgloss.NewStyle().Width(60).Foreground(lipgloss.Color("47"))
)

type model struct {
	w, h int // weight and height of screen

	activeView tea.Model

	hhclient *apiclient.ApiClient
	history  *views.History
}

func (m model) Init() tea.Cmd { return nil }

func (m model) View() string {
	return m.activeView.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case events.Message:
		m.history.Save(m.activeView)
		m.activeView = msg
	case events.Apply:
		var formText strings.Builder
		fmt.Fprint(&formText, "Откликнуться на вакансию?\n\n")
		fmt.Fprintln(&formText, views.RenderVacancyTitle(msg.VacancyNumber, msg.Title, msg.Salary, msg.Archived, msg.ResponseRequired))
		fmt.Fprintln(&formText, views.CompanyStyle.Render(msg.Employer))

		m.history.Save(m.activeView)
		m.activeView = forms.NewApplyForm(msg.VacancyID, formText.String(), m.hhclient)
	case events.ShowVacancy:
		var err error
		m.history.Save(m.activeView)
		if m.activeView, err = views.NewShowVacancy(
			msg.Name, msg.Salary, msg.Skills, msg.WorkFormat, msg.Desc, msg.Employer, m.w, m.h,
		); err != nil {
			m.activeView = events.NewMessage(events.ErrorMessage,
				fmt.Errorf("Show vacancy failure: %w", err).Error())
			return m, nil
		}
	case events.Source:
		var err error
		m.history.Save(m.activeView)
		if m.activeView, err = views.NewSource(msg.Vacancy, m.w, m.h); err != nil {
			m.activeView = events.NewMessage(events.ErrorMessage,
				fmt.Errorf("Show source(JSON) of vacancy data failure: %w", err).Error())
			return m, nil
		}
	case events.QuitFromInnerView:
		m.activeView = m.history.Back()
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	}

	var cmd tea.Cmd
	m.activeView, cmd = m.activeView.Update(msg)
	return m, cmd
}

func initialModel(
	client *apiclient.ApiClient,
	vacancies []dto.Vacancy,
	dict dto.Dictionary,
) model {
	hist := views.NewHistory()
	lv := views.NewListVacancies(vacancies, dict, client)
	hist.Save(lv)

	return model{
		hhclient:   client,
		activeView: lv,
		history:    hist,
	}
}
