package state

import (
	"time"

	"gitlab.com/1buran/hhbot/internal/infrastructure/apiclient/dto"
)

type Note struct {
	Text string
	Date time.Time
}

func NewNote(txt string) Note { return Note{Date: time.Now(), Text: txt} }

// The state is exported for easy save to & load from cache.
type State struct {
	BlacklistedVacancy map[string]Note
	BlacklistedCompany map[string]Note

	VacancyNotes map[string][]Note
	CompanyNotes map[string][]Note

	Vacancies map[string]dto.Vacancy

	Resumes  []dto.Resume
	ResumeID string // used in apply to vacancy
}

// Facade hides a state, all state mutations should be trough facade.
type Facade struct {
	state *State
	dirty bool
}

func (f Facade) ExportState() *State { return f.state }
func (f Facade) ResumeID() string    { return f.state.ResumeID }

func (f Facade) IsCompanyBlacklisted(companyID string) (*Note, bool) {
	meta, ok := f.state.BlacklistedCompany[companyID]
	return &meta, ok
}

func (f Facade) IsVacancyBlacklisted(vacancyID string) (*Note, bool) {
	meta, ok := f.state.BlacklistedVacancy[vacancyID]
	return &meta, ok
}

func (f *Facade) BlacklistCompany(companyID string, note Note) {
	f.dirty = true
	f.state.BlacklistedCompany[companyID] = note
}

func (f *Facade) BlacklistVacancy(vacancyID string, note Note) {
	f.dirty = true
	f.state.BlacklistedVacancy[vacancyID] = note
}

func (f *Facade) UnblacklistCompany(companyID string, note Note) {
	f.dirty = true
	delete(f.state.BlacklistedCompany, companyID)
}

func (f *Facade) UnblacklistVacancy(vacancyID string, note Note) {
	f.dirty = true
	delete(f.state.BlacklistedVacancy, vacancyID)
}

func (f *Facade) SetDefaultResume(resumeID string) {
	f.dirty = true
	f.state.ResumeID = resumeID
}

func (f *Facade) Init(state *State) { f.state = state }
func (f *Facade) SetClean()         { f.dirty = false }
func (f *Facade) IsDirty() bool     { return f.dirty }

func NewState() *Facade {
	return &Facade{
		state: &State{
			BlacklistedVacancy: make(map[string]Note),
			BlacklistedCompany: make(map[string]Note),
			VacancyNotes:       make(map[string][]Note),
			CompanyNotes:       make(map[string][]Note),
			Vacancies:          make(map[string]dto.Vacancy),
		},
	}
}
