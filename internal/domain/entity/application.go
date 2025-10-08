package entity

import "time"

// An Application holds information about applicant and vacancy for which he/she has applied.
type Application struct {
	id      string
	status  string // invite to interview, discard etc
	created time.Time
	person  Person
	vacancy Vacancy
}

func (a Application) GetID() string { return a.id }
