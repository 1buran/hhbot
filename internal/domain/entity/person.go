package entity

type Person struct {
	id   string
	name string
	url  string
}

func (p Person) GetID() string { return p.id }
