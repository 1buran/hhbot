package dto

import (
	"fmt"
	"strings"
	"time"
)

type Salary struct {
	Currency string
	From, To int
	Gross    bool
}

func (s Salary) String() string {
	var out []string
	if s.From > 0 {
		out = append(out, fmt.Sprintf("от %d", s.From))
	}
	if s.To > 0 {
		out = append(out, fmt.Sprintf("до %d", s.To))
	}
	if len(out) > 0 {
		out = append(out, s.Currency)
	}
	return strings.Join(out, " ")
}

type Counters struct{ Responses, Total_responses int }

type Snippet struct{ Requirement, Responsibility string }
type Skill struct{ Name string }

type Vacancy struct {
	Id                  string
	Name                string
	Code                string
	Description         string
	Alternate_url       string
	Apply_alternate_url string
	Response_url        string
	Url                 string

	Relations  []string // favorited, got_response, got_invitation, got_rejection, blacklisted
	Key_skills []Skill

	// Do not forget enable json v2 for get new features like time format: GOEXPERIMENT=jsonv2
	Created_at         time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
	Published_at       time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
	Initial_created_at time.Time `json:",format:'2006-01-02T15:04:05-0700'"`

	Employer Employer

	Professional_roles DictionaryItems

	Salary_range Salary

	Type                  DictionaryItem
	Employment_form       DictionaryItem
	Experience            DictionaryItem
	Work_format           DictionaryItems
	Work_schedule_by_days DictionaryItems
	Working_hours         DictionaryItems

	Counters Counters

	Snippet Snippet

	Premium  bool
	Archived bool
	Approved bool

	Show_contacts            bool
	Allow_messages           bool
	Closed_for_applicants    bool
	Response_letter_required bool
}
