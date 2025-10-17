package views

const (
	ListVacancies = iota // views, screens with different content type
	VacancyContent
	ApplyToVacancy
	VacancySource
)

type RoutingState struct {
	current, prev, home uint8
}

func (d *RoutingState) Update(v uint8) {
	d.prev = d.current
	d.current = v
}

func (d *RoutingState) Back() {
	d.prev = d.home
	d.current = d.prev
}

func (d *RoutingState) Home() {
	d.prev = d.current
	d.current = d.home
}

func (d RoutingState) Current() uint8 { return d.current }

func NewRoutingState() *RoutingState { return &RoutingState{} }
