package entity

type Company struct {
	id      string
	name    string
	url     string
	trusted bool // company should passed moderation
}

func (c Company) GetID() string { return c.id }
