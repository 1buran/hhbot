package main

import (
	"fmt"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	_ "github.com/joho/godotenv/autoload"
	"github.com/muesli/termenv"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
	"strings"
)

var (
	glamourRenderer, _ = glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithColorProfile(termenv.TrueColor),
	)
)

type model struct {
	vacancies []dto.Vacancy // load vacancies from hh.ru
	cursor    int           // active vacancy
	view      uint8         // alternate views

	content []string // splitted by lines for scroll
	scroll  int      // vertical scroll position, beginning(first line) of visible content part
	w, h    int      // weight and height of screen

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
		lines := len(m.content)
		if lines <= m.h { // all content fits the one screen without scroll
			return strings.Join(m.content, "\n")
		}
		return strings.Join(m.content[m.scroll:m.scroll+m.h-1], "\n")
	}
	return renderVacancy(m.cursor, m.vacancies[m.cursor], m.hhdict)
}

type viewVacancy struct{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case viewVacancy, error:
		m.view = 1
		return m, nil
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
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
			switch m.view {
			case 0:
				if m.cursor > 0 {
					m.cursor--
				}
			case 1:
				if m.scroll > 0 {
					m.scroll--
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			switch m.view {
			case 0:
				if m.cursor < len(m.vacancies)-1 {
					m.cursor++
				}
			case 1:
				if len(m.content)-m.scroll > m.h {
					m.scroll++
				}
			}
		// Enter to vacancy
		case "enter", " ":
			m.scroll = 0
			m.content = []string{}
			vacancyID := m.vacancies[m.cursor].Id
			vac, err := m.hhclient.GetVacancy(vacancyID)
			if err != nil {
				m.content = []string{redflagStyle.Render(err.Error())}
				return m, func() tea.Msg {
					return fmt.Errorf("GetVacancy failure: %w", err)
				}
			}
			md, err := html2md.ConvertString(vac.Description)
			if err != nil {
				m.content = []string{redflagStyle.Render(err.Error())}
				return m, func() tea.Msg {
					return fmt.Errorf("html2md.ConvertString failure: %w", err)
				}
			}

			renderedContent, err := glamourRenderer.Render(md)
			if err != nil {
				m.content = []string{redflagStyle.Render(err.Error())}
				return m, func() tea.Msg {
					return fmt.Errorf("glamour.Render failure: %w", err)
				}
			}
			m.content = strings.Split(renderedContent, "\n")
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
