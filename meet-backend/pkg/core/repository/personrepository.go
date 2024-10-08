package repository

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
)

type PersonRepository interface {
	FindById(uuid string) (*domain.Person, error)

	FindByEmail(email string) (*domain.Person, error)

	FindAll() ([]domain.Person, error)

	FilterPaginated(filters domain.PersonFilter) ([]domain.Person, error)

	Persist(person domain.Person) (*domain.Person, error)

	Update(person domain.Person) (*domain.Person, error)

	Delete(uuid string) error
}
