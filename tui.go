package main

import (
	"encoding/json"
	"fmt"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/joho/godotenv/autoload"
	"github.com/muesli/termenv"

	"strings"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

const (
	ListVacancies = iota // views, screens with different content type
	VacancyContent
	ApplyToVacancy
	ErrorMessage
	Information
	VacancySource
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

	applyForm Form

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
		// todo: make it like a modal window
		return informationStyle.Render(m.msg)
	case ErrorMessage:
		return redflagStyle.Render(m.msg)
	case VacancyContent, VacancySource:
		lines := len(m.content)
		if lines <= m.h { // all content fits the one screen without scroll
			return strings.Join(m.content, "\n")
		}
		return strings.Join(m.content[m.scroll:m.scroll+m.h-1], "\n")
	case ApplyToVacancy:
		return m.applyForm.View()
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
		// todo: refactoring exit from textarea focused input
		// currently it is on escape
		if m.applyForm != nil {
			if m.applyForm.Focused() && msg.String() != "esc" {
				break
			}
		}

		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			if m.view == VacancyContent { // quit from preview of vacancy
				m.view = ListVacancies
				return m, nil
			}
			return m, tea.Quit
		case "tab":
			switch m.view {
			case ApplyToVacancy:
				return m, m.applyForm.Tab()
			}
		case "a":
			// todo: refresh vacancy info before apply: data gathered from search
			// is not full - a lot of fields are not filled, but will be filled on
			// GET /vacancies/{id}
			if m.vacancies[m.cursor].Test.Required {
				m.view = ErrorMessage
				m.msg = fmt.Sprintf("%s\n%s: %s",
					"Необходимо пройти опрос/тест работодателя",
					"переходите по ссылке и подавайтесь в ручную",
					m.vacancies[m.cursor].Alternate_url)
				return m, nil
			}
			v := m.vacancies[m.cursor]
			m.applyForm = NewDialogYesNo("Откликнуться на вакансию?\n" +
				renderVacancyTitle(m.cursor, v) + "\n" +
				companyStyle.Render(v.Employer.Name),
			)
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
					"b96ed6a4ff0f1ea45a0039ed1f3153446c5763",
					m.applyForm.Message(),
				)
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
					m.msg = lipgloss.NewStyle().Width(60).Render(res.Error)
				default:
					m.view = ErrorMessage
					m.msg = fmt.Sprintf("%+v", res)
				}
			}
		case "esc":
			if m.applyForm != nil && m.applyForm.Focused() {
				return m, m.applyForm.Tab()
			}
			return m, goBack(m)

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			switch m.view {
			case ListVacancies:
				if m.cursor > 0 {
					m.cursor--
				}
			case VacancyContent, VacancySource:
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
			case VacancyContent, VacancySource:
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

			m.vacancies[m.cursor] = vac // update memory object with additional data

			title := fmt.Sprintf("**%s** / `%s` / _%s_\n\n---\n\n",
				vac.Name, vac.Salary_range.String(), vac.Employer.Name)
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

			renderedContent, err := glamourRenderer.Render(title + skills + md)
			if err != nil {
				m.view = ErrorMessage
				m.msg = fmt.Errorf("glamour.Render failure: %w", err).Error()
				return m, nil
			}

			m.view = VacancyContent
			m.content = strings.Split(renderedContent, "\n")

		case "s":
			b, err := json.MarshalIndent(m.vacancies[m.cursor], "", "    ")
			if err != nil {
				m.view = ErrorMessage
				m.msg = err.Error()
			}
			content := fmt.Sprintf("**Исходное содержание объекта вакансия:**\n\n ```json\n\n%s\n\n```", b)
			renderedContent, err := glamourRenderer.Render(content)
			if err != nil {
				m.view = ErrorMessage
				m.msg = fmt.Errorf("glamour.Render failure: %w", err).Error()
				return m, nil
			}

			m.scroll = 0
			m.view = VacancySource
			m.prev = ListVacancies // m.contnet  will be overwritten so it is need to load vacancy again
			m.content = strings.Split(renderedContent, "\n")
		}
	}

	var cmds []tea.Cmd
	if m.applyForm != nil {
		cmds = append(cmds, m.applyForm.Update(msg))
	}
	return m, tea.Batch(cmds...)
}

func initialModel(
	client *apiclient.ApiClient,
	vacancies []dto.Vacancy,
	dict dto.Dictionary,
) model {
	return model{vacancies: vacancies, hhdict: dict, hhclient: client}
}
