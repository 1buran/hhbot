package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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

	activeView tea.Model
	statusbar  tea.Model

	hhclient *apiclient.ApiClient
	hhdict   dto.Dictionary
	state    *state.Facade

	history *views.History
}

func (m model) Init() tea.Cmd { return m.activeView.Init() }

func (m model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, m.activeView.View(), "", m.statusbar.View())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.statusbar, _ = m.statusbar.Update(msg) // statusbar should not fires any commands

	switch msg := msg.(type) {
	case actions.Blacklisted:
		if msg.CompanyID != "" {
			// todo: remove hardcode, show user dialog for input text of note
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
		m.history.Save(m.activeView, m.statusbar)
		m.activeView = msg
		lvl := events.StatusBarNotifyLevelDefault
		if msg.Level() == events.ErrorMessage {
			lvl = events.StatusBarNotifyLevelAlert
		}
		m.statusbar = views.NewStatusBarWithContext(lvl, msg.Status(), m.w)
	case events.NewSearch:
		m.history.Save(m.activeView, m.statusbar)
		m.activeView = forms.NewInput(m.hhclient)
		m.statusbar = views.NewStatusBar()
	case events.SearchResults:
		m.history.Save(m.activeView, m.statusbar)
		m.activeView = views.NewListVacancies(msg.Vacancies, m.hhdict, m.hhclient, m.state)
		m.statusbar = views.NewStatusBarWithContext(events.StatusBarNotifyLevelDefault,
			fmt.Sprintf("Найдено %d вакансий", len(msg.Vacancies)), m.w)
	case events.Apply:
		var formText strings.Builder
		fmt.Fprint(&formText, "Откликнуться на вакансию?\n\n")
		fmt.Fprintln(&formText, views.RenderVacancyTitle(msg.VacancyNumber, msg.Title, msg.Salary, msg.Archived, msg.ResponseRequired))
		fmt.Fprintln(&formText, styles.Company.Render(msg.Employer))

		m.history.Save(m.activeView, m.statusbar)
		m.activeView = forms.NewApplyForm(msg.VacancyID, m.state.ResumeID(), formText.String(), m.hhclient)
	case events.ShowVacancy:
		var err error
		m.history.Save(m.activeView, m.statusbar)
		if m.activeView, err = views.NewShowVacancy(
			msg.Name, msg.Salary, msg.Skills, msg.WorkFormat, msg.Desc, msg.Employer, m.w, m.h,
		); err != nil {
			m.activeView = events.NewMessage(events.ErrorMessage,
				fmt.Errorf("show vacancy failure: %w", err).Error(),
				"Ошибка просмотра вакансии")
			return m, nil
		}
	case events.Source:
		var err error
		m.history.Save(m.activeView, m.statusbar)
		if m.activeView, err = views.NewSource(msg.Vacancy, m.w, m.h); err != nil {
			m.activeView = events.NewMessage(events.ErrorMessage,
				fmt.Errorf("show source(JSON) of vacancy data failure: %w", err).Error(),
				"Ошибка просмотра исходного JSON")
			return m, nil
		}
	case events.QuitFromInnerView:
		views := m.history.Back()
		if len(views) == 2 {
			m.activeView, m.statusbar = views[0], views[1]
		}
	case events.QuitFromResumeLoader:
		m.activeView = forms.NewInput(m.hhclient)
		m.statusbar = views.NewStatusBar()
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	}

	m.activeView, cmd = m.activeView.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func InitialModel(
	activeView tea.Model,
	client *apiclient.ApiClient,
	state *state.Facade,
	dict dto.Dictionary,
) model {
	hist := views.NewHistory()
	hist.Save(activeView)

	sbar := views.NewStatusBar()
	return model{
		hhclient:   client,
		hhdict:     dict,
		activeView: activeView,
		statusbar:  sbar,
		history:    hist,
		state:      state,
	}
}
