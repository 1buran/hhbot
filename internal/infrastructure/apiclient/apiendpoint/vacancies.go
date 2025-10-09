package apiendpoint

import (
	"net/url"
	"strconv"
)

type getvacancies struct {
	text                    string
	schedule                string
	experience              string
	employment              string
	employer_id             string
	excluded_employer_id    string
	professional_role       string
	currency                string
	label                   string
	date_from, date_to      string
	part_time               string
	salary                  int
	no_magic                bool // disable automatic conversion of original textual query to query params
	only_with_salary        bool // show vacancies only with specified salary range, set to true
	responses_count_enabled bool // include counters to response, set to true

	Pagination
}

func (v *getvacancies) SetPage(p int)           { v.Pagination.Page = p }
func (v getvacancies) Url(params ...any) string { return "/vacancies" }
func (v getvacancies) Payload() url.Values {
	payload := url.Values{}

	if v.text != "" {
		payload.Add("text", v.text)
	}
	if v.schedule != "" {
		payload.Add("schedule", v.schedule)
	}
	if v.experience != "" {
		payload.Add("experience", v.experience)
	}
	if v.employment != "" {
		payload.Add("employment", v.employment)
	}
	if v.employer_id != "" {
		payload.Add("employer_id", v.employer_id)
	}
	if v.excluded_employer_id != "" {
		payload.Add("excluded_employer_id", v.excluded_employer_id)
	}
	if v.professional_role != "" {
		payload.Add("professional_role", v.professional_role)
	}
	if v.currency != "" {
		payload.Add("currency", v.currency)
	}
	if v.label != "" {
		payload.Add("label", v.label)
	}
	if v.date_from != "" {
		payload.Add("date_from", v.date_from)
	}
	if v.date_to != "" {
		payload.Add("date_to", v.date_to)
	}
	if v.part_time != "" {
		payload.Add("part_time", v.part_time)
	}

	if v.salary > 0 {
		payload.Add("salary", strconv.Itoa(v.salary))
	}
	if v.Pagination.Page > 0 {
		payload.Add("page", strconv.Itoa(v.Pagination.Page))
	}
	if v.Pagination.PerPage > 0 {
		payload.Add("per_page", strconv.Itoa(v.Pagination.PerPage))
	}

	if v.no_magic {
		payload.Add("no_magic", "true")
	}
	if v.only_with_salary {
		payload.Add("only_with_salary", "true")
	}
	if v.responses_count_enabled {
		payload.Add("responses_count_enabled", "true")
	}

	return payload
}

func NewGetVacancies(text string, salary int) ApiEndpoint {
	return &getvacancies{
		// defaults:
		only_with_salary:        true,
		responses_count_enabled: true,
		currency:                "RUR",

		// query params:
		text:   text,
		salary: salary,

		// pagination settings:
		Pagination: Pagination{PerPage: 20},
	}
}
