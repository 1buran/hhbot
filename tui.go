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

const (
	ListVacancies = iota // views, screens with different content type
	VacancyContent
	ApplyToVacancy
	ErrorMessage
	Information
)

var (
	glamourRenderer, _ = glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithColorProfile(termenv.TrueColor),
	)
)

type model struct {
	vacancies  []dto.Vacancy // load vacancies from hh.ru
	cursor     int           // active vacancy
	view, prev uint8         // alternate views

	msg     string
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
	case Information:
		// todo make it like a modal window
		return informationStyle.Render(m.msg)
	case ErrorMessage:
		return redflagStyle.Render(m.msg)
	case VacancyContent:
		lines := len(m.content)
		if lines <= m.h { // all content fits the one screen without scroll
			return strings.Join(m.content, "\n")
		}
		return strings.Join(m.content[m.scroll:m.scroll+m.h-1], "\n")
	case ApplyToVacancy:
		v := m.vacancies[m.cursor]
		return DialogYesNo("Откликнуться на вакансию?\n" +
			renderVacancyTitle(m.cursor, v) + "\n" +
			companyStyle.Render(v.Employer.Name),
		)
	}

	// default: show list of vacancies
	return renderVacancy(m.cursor, m.vacancies[m.cursor], m.hhdict)
}

// Route switch between screens.
type Routing struct {
	view uint8
}

// Go back to previous screen.
func goBack(m model) func() tea.Msg {
	return func() tea.Msg {
		return Routing{view: m.prev}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Routing:
		m.prev = m.view
		m.view = msg.view
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			if m.view == VacancyContent { // quit from preview of vacancy
				m.view = ListVacancies
				return m, nil
			}
			return m, tea.Quit
		case "a":
			m.prev = m.view
			m.view = ApplyToVacancy
		case "n":
			if m.view == ApplyToVacancy {
				return m, goBack(m)
			}
		case "y":
			if m.view == ApplyToVacancy {
				vacancyID := m.vacancies[m.cursor].Id
				res, err := m.hhclient.Apply(vacancyID,
					"b96ed6a4ff0f1ea45a0039ed1f3153446c5763", "Интересно!")
				if err != nil {
					m.view = ErrorMessage
					m.msg = err.Error()
					return m, nil
				}
				switch res.Code {
				case 201:
					return m, goBack(m)
				case 303:
					m.view = Information
					m.msg = fmt.Sprintf("Перейдите на: %s", res.Headers.Get("Location"))
				case 400, 403:
					m.view = ErrorMessage
					m.msg = res.Error
				default:
					m.view = ErrorMessage
					m.msg = fmt.Sprintf("%+v", res)
				}
			}
		case "esc":
			return m, goBack(m)

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			switch m.view {
			case ListVacancies:
				if m.cursor > 0 {
					m.cursor--
				}
			case VacancyContent:
				if m.scroll > 0 {
					m.scroll--
				}
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			switch m.view {
			case ListVacancies:
				if m.cursor < len(m.vacancies)-1 {
					m.cursor++
				}
			case VacancyContent:
				if len(m.content)-m.scroll > m.h {
					m.scroll++
				}
			}
		// Enter to vacancy
		case "enter", " ":
			m.scroll = 0
			m.content = []string{}
			m.prev = m.view
			vacancyID := m.vacancies[m.cursor].Id
			vac, err := m.hhclient.GetVacancy(vacancyID)
			if err != nil {
				m.view = ErrorMessage
				m.msg = fmt.Errorf("GetVacancy failure: %w", err).Error()
				return m, nil
			}

			var skills string
			if vac.Key_skills.Length() > 0 {
				skills = fmt.Sprintf("**Skills**: %s\n\n", vac.Key_skills.String())
			}

			md, err := html2md.ConvertString(vac.Description)
			if err != nil {
				m.view = ErrorMessage
				m.msg = fmt.Errorf("html2md.ConvertString failure: %w", err).Error()
				return m, nil
			}

			renderedContent, err := glamourRenderer.Render(skills + md)
			if err != nil {
				m.view = ErrorMessage
				m.msg = fmt.Errorf("glamour.Render failure: %w", err).Error()
				return m, nil
			}

			m.view = VacancyContent
			m.content = strings.Split(renderedContent, "\n")
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
