package apiendpoint

import (
	"net/url"
)

type getvacancy struct {
	vacancy_id string
}

func (v getvacancy) SetPage(p int)            {}
func (v getvacancy) Url(params ...any) string { return "/vacancies/" + v.vacancy_id }
func (v getvacancy) Payload() url.Values      { return url.Values{} }

func NewGetVacancy(vacancyID string) ApiEndpoint {
	return &getvacancy{vacancy_id: vacancyID}
}
