package dto

import (
	"time"
)

type ResumeSalary struct {
	Amount   int
	Currency string
}

func (s ResumeSalary) String() string {
	return printer.Sprintf("%d %s", s.Amount, s.Currency)
}

type ResumeExperience struct {
	Months int
}

type Resume struct {
	Id, Real_id string

	Title            string
	Salary           ResumeSalary
	Total_experience ResumeExperience

	First_name, Last_name, Middle_name string

	Age                    int
	New_views, Total_views int

	Alternate_url, Url string

	Can_view_full_info, Blocked, Finished, Can_publish_or_update bool

	Created_at      time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
	Updated_at      time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
	Next_publish_at time.Time `json:",format:'2006-01-02T15:04:05-0700'"`
}
