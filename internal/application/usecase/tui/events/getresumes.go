package events

import (
	tea "github.com/charmbracelet/bubbletea"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type ResumesLoaded struct {
	Resumes []dto.Resume
}

func NewResumesLoaded(resumes []dto.Resume) ResumesLoaded { return ResumesLoaded{Resumes: resumes} }

type QuitFromResumeLoader struct{}

func NewQuitFromResumeLoader() tea.Cmd {
	return func() tea.Msg {
		return QuitFromResumeLoader{}
	}
}
