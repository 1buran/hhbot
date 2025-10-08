package repository

import "gitlab.com/1buran/hhbot/internal/domain/entity"

type Entities interface {
	GetID() string

	entity.Vacancy | entity.Company | entity.Person | entity.Application
}

type Repository[T Entities] interface {
	Add(v T)
	List() []T
}
