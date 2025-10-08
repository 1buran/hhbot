package apiclient

import "time"

type DictionaryItemDTO struct{ Id, Name string }

type EmployerDTO struct {
	Id            string
	Name          string
	Trusted       bool
	Url           string
	Alternate_url string
}

type SalaryDTO struct {
	Currency string
	From, To int
	Gross    bool
}

type CountersDTO struct{ Responses, Total_responses int }

type SnippetDTO struct{ Requirement, Responsibility string }

type VacancyDTO struct {
	Id                  string
	Name                string
	Alternate_url       string
	Apply_alternate_url string
	Response_url        string
	Url                 string

	Relations []string // favorited, got_response, got_invitation, got_rejection, blacklisted

	// Do not forget enable json v2 for get new features like time format: GOEXPERIMENT=jsonv2
	Created_at   time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
	Published_at time.Time `json:",format:'2006-01-02T15:04:05-0700'"`

	Employer EmployerDTO

	Professional_roles []DictionaryItemDTO

	Salary_range SalaryDTO

	Type                  DictionaryItemDTO
	Employment_form       DictionaryItemDTO
	Experience            DictionaryItemDTO
	Work_format           []DictionaryItemDTO
	Work_schedule_by_days []DictionaryItemDTO
	Working_hours         []DictionaryItemDTO

	Counters CountersDTO

	Snippet SnippetDTO

	Premium  bool
	Archived bool
}

type ResponseType interface {
	VacancyDTO
}

type ResponseItemsDTO[T ResponseType] struct {
	Found, Per_page, Pages, Page int
	Items                        []T
}
