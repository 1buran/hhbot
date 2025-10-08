package entity

import "time"

// A Vacancy holds the information about the published vacancy.
type Vacancy struct {
	id           string
	name         string
	description  string
	url          string
	applyurl     string
	worktype     string // full, partial etc
	workformat   string // remote, hybrid etc
	workschedule string // 5/2, 7/0 etc
	experience   string // e.g. between 1 to 3 years
	salary       string
	skills       []string
	published    time.Time

	approved bool // vacancy should passed moderation
	archived bool // it maybe archived
	closed   bool // or closed for applicants

	company Company
}

func (v Vacancy) GetID() string { return v.id }
