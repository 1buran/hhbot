package views

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	_ "github.com/joho/godotenv/autoload"
	"github.com/muesli/termenv"

	"gitlab.com/1buran/hhbot/internal/application/usecase/tui/styles"
	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

var (
	glamourRenderer, _ = glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithWordWrap(100),
		glamour.WithColorProfile(termenv.TrueColor),
	)
)

func renderExperience(exp dto.DictionaryItem) (experience string) {
	switch exp.Id {
	case "moreThan6":
		experience = styles.Redflag.Render(exp.Name)
	case "between1And3", "noExperience":
		experience = styles.Greenlight.Render(exp.Name)
	default:
		experience = exp.Name
	}
	return
}

func renderRelations(relations []string, dict dto.Dictionary) string {
	relationsDict, ok := dict.Get("vacancy_relation")
	if !ok {
		fmt.Printf("Warning: vacancy_relation not found! dict: %+v\n", dict)
	}

	var rel []string
	for _, k := range relations {
		v := relationsDict.Get(k)
		if v != "" {
			switch k {
			case "favorite":
				v = styles.Favorite.Render(v)
			case "got_response":
				v = styles.Applied.Render(v)
			case "got_rejection":
				v = styles.Redflag.Render(v)
			case "got_invitation":
				v = styles.Invited.Render(v)
			case "blacklisted":
				v = styles.Blacklisted.Render(v)
			case "got_question":
				v = styles.Question.Render(v)
			}
			rel = append(rel, v)
		}
	}
	return strings.Join(rel, ", ")
}

func RenderVacancyTitle(i int, title, salary string, archived, respRequired bool) string {
	title = styles.Title.Render(title)
	if archived {
		title += styles.Redflag.Render(" ВНИМАНИЕ! Вакансия в архиве!")
	}
	if respRequired {
		title += styles.Redflag.Render(" ВНИМАНИЕ! Требуется сопроводительное письмо при отклике!")
	}
	return fmt.Sprintf("%d. %s / %s\n", i+1, title,
		styles.Salary.Render(salary))
}

func renderVacancy(i int, v dto.Vacancy, dict dto.Dictionary) string {
	renderedRelations := renderRelations(v.Relations, dict)
	experience := renderExperience(v.Experience)

	var buf strings.Builder
	fmt.Fprint(&buf, RenderVacancyTitle(i, v.Name, v.Salary_range.String(), v.Archived, v.Response_letter_required))
	fmt.Fprintln(&buf,
		styles.VacancyCard.Render(
			fmt.Sprintf(
				"Компания: %s\nОпыт: %s\nСвязь: %s\nФормат работы: %s\n"+
					"Откликнулось: %d / %d\nОпубликовано: %s\nСоздано: %s\n%s\n\n",
				styles.Company.Render(v.Employer.Name),
				experience,
				renderedRelations,
				v.Work_format.String(),
				v.Counters.Responses, v.Counters.Total_responses,
				v.Published_at.Format(time.DateOnly),
				v.Created_at.Format(time.DateOnly),
				v.Alternate_url,
			),
		),
	)
	return buf.String()
}
