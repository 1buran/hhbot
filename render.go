package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	_ "github.com/joho/godotenv/autoload"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

var (
	vacancyCard   = lipgloss.NewStyle().PaddingLeft(4)
	redflagStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6863"))
	appliedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#80EF80"))
	invitedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	questionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00F0FF"))
	favoriteStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF7518"))

	blacklistedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFEF00"))
	greenlightStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("121"))

	salaryStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#dddddd")).Bold(true)
	companyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Italic(true)

	informationStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("47"))
)

func renderExperience(exp dto.DictionaryItem) (experience string) {
	switch exp.Id {
	case "moreThan6":
		experience = redflagStyle.Render(exp.Name)
	case "between1And3", "noExperience":
		experience = greenlightStyle.Render(exp.Name)
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
				v = favoriteStyle.Render(v)
			case "got_response":
				v = appliedStyle.Render(v)
			case "got_rejection":
				v = redflagStyle.Render(v)
			case "got_invitation":
				v = invitedStyle.Render(v)
			case "blacklisted":
				v = blacklistedStyle.Render(v)
			case "got_question":
				v = questionStyle.Render(v)
			}
			rel = append(rel, v)
		}
	}
	return strings.Join(rel, ", ")
}

func renderVacancyTitle(i int, v dto.Vacancy) string {
	title := titleStyle.Render(v.Name)
	if v.Archived {
		title += redflagStyle.Render(" ВНИМАНИЕ! Вакансия в архиве!")
	}
	if v.Response_letter_required {
		title += redflagStyle.Render(" ВНИМАНИЕ! Требуется сопроводительное письмо при отклике!")
	}
	return fmt.Sprintf("%d. %s / %s\n", i+1, title,
		salaryStyle.Render(v.Salary_range.String()))
}

func renderVacancy(i int, v dto.Vacancy, dict dto.Dictionary) string {
	renderedRelations := renderRelations(v.Relations, dict)
	experience := renderExperience(v.Experience)

	var buf strings.Builder
	fmt.Fprint(&buf, renderVacancyTitle(i, v))
	fmt.Fprintln(&buf,
		vacancyCard.Render(
			fmt.Sprintf(
				"Компания: %s\nОпыт: %s\nСвязь: %s\nФормат работы: %s\n"+
					"Откликнулось: %d / %d\nОпубликовано: %s\nСоздано: %s\n%s\n\n",
				companyStyle.Render(v.Employer.Name),
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
