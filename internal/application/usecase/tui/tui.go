package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/actions"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/forms"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/state"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/views"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type model struct {
	w, h int // weight and height of screen

	resumeID string

	activeView tea.Model

	hhclient *apiclient.ApiClient
	hhdict   dto.Dictionary
	history  *views.History
	state    *state.Facade
}

func (m model) Init() tea.Cmd { return nil }

func (m model) View() string {
	return m.activeView.View()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case actions.Blacklisted:
		if msg.CompanyID != "" {
			m.state.BlacklistCompany(msg.CompanyID, state.NewNote(
				"Компания добавлена в чёрный список"))
		}
		if msg.VacancyID != "" {
			m.state.BlacklistVacancy(msg.VacancyID, state.NewNote(
				"Вакансия добавлена в чёрный список"))
		}
	case actions.Unblacklisted:
		if msg.CompanyID != "" {
			m.state.UnblacklistCompany(msg.CompanyID, state.NewNote(
				"Компания удалена из чёрного спискa"))
		}
		if msg.VacancyID != "" {
			m.state.UnblacklistVacancy(msg.VacancyID, state.NewNote(
				"Вакансия удалена из чёрного списка"))
		}
	case events.Message:
		m.history.Save(m.activeView)
		m.activeView = msg
	case events.NewSearch:
		m.history.Save(m.activeView)
		m.activeView = forms.NewInput(m.hhclient)
	case events.SearchResults:
		m.history.Save(m.activeView)
		m.activeView = views.NewListVacancies(msg.Vacancies, m.hhdict, m.hhclient, m.state)
	case events.Apply:
		var formText strings.Builder
		fmt.Fprint(&formText, "Откликнуться на вакансию?\n\n")
		fmt.Fprintln(&formText, views.RenderVacancyTitle(msg.VacancyNumber, msg.Title, msg.Salary, msg.Archived, msg.ResponseRequired))
		fmt.Fprintln(&formText, styles.Company.Render(msg.Employer))

		m.history.Save(m.activeView)
		m.activeView = forms.NewApplyForm(msg.VacancyID, m.resumeID, formText.String(), m.hhclient)
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

func InitialModel(
	activeView tea.Model,
	client *apiclient.ApiClient,
	state *state.Facade,
	dict dto.Dictionary,
	resumeID string,
) model {
	hist := views.NewHistory()
	hist.Save(activeView)

	return model{
		hhclient:   client,
		hhdict:     dict,
		activeView: activeView,
		history:    hist,
		state:      state,
		resumeID:   resumeID,
	}
}
