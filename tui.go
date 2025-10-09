package main

import (
	"fmt"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type model struct {
	vacancies []dto.Vacancy // load vacancies from hh.ru
	cursor    int           // active vacancy
	view      uint8         // alternate views

	content string

	hhdict   dto.Dictionary
	hhclient *apiclient.ApiClient
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) View() string {
	switch m.view {
	case 1:
		return m.content
	}
	return renderVacancy(m.cursor, m.vacancies[m.cursor], m.hhdict)
}

type viewVacancy struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewVacancy, error:
		m.view = 1
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			if m.view == 1 { // quit from preview of vacancy
				m.view = 0
				return m, nil
			}
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.vacancies)-1 {
				m.cursor++
			}

		// Enter to vacancy
		case "enter", " ":
			m.content = ""
			vacancyID := m.vacancies[m.cursor].Id
			vac, err := m.hhclient.GetVacancy(vacancyID)
			if err != nil {
				m.content = redflagStyle.Render(err.Error())
				return m, func() tea.Msg {
					return fmt.Errorf("GetVacancy failure: %w", err)
				}
			}
			md, err := html2md.ConvertString(vac.Description)
			if err != nil {
				m.content = redflagStyle.Render(err.Error())
				return m, func() tea.Msg {
					return fmt.Errorf("html2md.ConvertString failure: %w", err)
				}
			}

			m.content, err = glamour.Render(md, "dark")
			if err != nil {
				m.content = redflagStyle.Render(err.Error())
				return m, func() tea.Msg {
					return fmt.Errorf("glamour.Render failure: %w", err)
				}
			}

			return m, func() tea.Msg {
				return viewVacancy{}
			}
		}
	}
	return m, nil
}

func initialModel(
	client *apiclient.ApiClient,
	vacancies []dto.Vacancy,
	dict dto.Dictionary,
) model {
	return model{vacancies: vacancies, hhdict: dict, hhclient: client}
}
