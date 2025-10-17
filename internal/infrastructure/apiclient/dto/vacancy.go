package dto

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/message"
)

var (
	printer = message.NewPrinter(message.MatchLanguage("ru"))
)

type Salary struct {
	Currency string
	From, To int
	Gross    bool
}

func (s Salary) String() string {
	var out []string
	if s.From > 0 {
		out = append(out, fmt.Sprint("от ", printer.Sprint(s.From)))
	}
	if s.To > 0 {
		out = append(out, fmt.Sprint("до ", printer.Sprint(s.To)))
	}
	if len(out) > 0 {
		out = append(out, s.Currency)
	}
	return strings.Join(out, " ")
}

type Counters struct{ Responses, Total_responses int }

type Snippet struct{ Requirement, Responsibility string }
type Skill struct{ Name string }

type Skills []Skill

func (s Skills) String() string {
	var tags []string
	for _, v := range s {
		tags = append(tags, v.Name)
	}
	return strings.Join(tags, ", ")
}

func (s Skills) Length() int { return len(s) }

type Test struct {
	Id       string
	Required bool
}

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
	Key_skills Skills

	// Do not forget enable json v2 for get new features like time format: GOEXPERIMENT=jsonv2
	Created_at         time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
	Published_at       time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
	Initial_created_at time.Time `json:",format:'2006-01-02T15:04:05-0700'"`

	Employer Employer

	Professional_roles DictionaryItems

	Salary_range Salary
	Test         Test

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
