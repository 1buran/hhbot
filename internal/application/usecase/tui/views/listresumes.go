package views

import (
	tea "github.com/charmbracelet/bubbletea"

	"fmt"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/events"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/state"
	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type listresumesModel struct {
	resumes []dto.Resume
	cursor  int

	state       *state.Facade
	loadResumes func() ([]dto.Resume, error)

	loaded bool
}

func (m listresumesModel) Init() tea.Cmd {
	return func() tea.Msg {
		resumes, err := m.loadResumes()
		if err != nil {
			return events.NewMessage(events.ErrorMessage,
				fmt.Errorf("get resumes failure: %w", err).Error(),
				"Ошибка загрузки резюме")
		}
		return events.NewResumesLoaded(resumes)
	}
}

func (m listresumesModel) View() string {
	if len(m.resumes) > 0 {
		r := m.resumes[m.cursor]
		return renderResume(r)
	}
	if !m.loaded {
		return styles.Information.Render("Загружаются резюме...")
	}
	return styles.Information.Render("Резюме не найдены!")
}

func (m listresumesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case events.ResumesLoaded:
		m.loaded = true
		m.resumes = msg.Resumes
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, events.NewQuitFromResumeLoader()
		case "enter":
			rid := m.resumes[m.cursor].Id
			m.state.SetDefaultResume(rid)
			return m, events.NewStatusBarNotifyCmd(
				"Резюме по умолчанию установлено, ResumeID:"+rid,
				events.StatusBarNotifyLevelDefault)
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.resumes)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func NewListResumes(client *apiclient.ApiClient, state *state.Facade) *listresumesModel {
	return &listresumesModel{state: state, loadResumes: client.GetResumes}
}
