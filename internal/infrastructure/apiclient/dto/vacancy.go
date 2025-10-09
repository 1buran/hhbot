package dto

import (
	"fmt"
	"time"
)

type Salary struct {
	Currency string
	From, To int
	Gross    bool
}

func (s Salary) String() string {
	var from, to any = "-", "-"
	if s.From > 0 {
		from = s.From
	}
	if s.To > 0 {
		from = s.To
	}
	return fmt.Sprintf("от: %v, до: %v %s", from, to, s.Currency)
}

type Counters struct{ Responses, Total_responses int }

type Snippet struct{ Requirement, Responsibility string }

type Vacancy struct {
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

	Employer Employer

	Professional_roles []DictionaryItem

	Salary_range Salary

	Type                  DictionaryItem
	Employment_form       DictionaryItem
	Experience            DictionaryItem
	Work_format           []DictionaryItem
	Work_schedule_by_days []DictionaryItem
	Working_hours         []DictionaryItem

	Counters Counters

	Snippet Snippet

	Premium  bool
	Archived bool
}
