package service

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/exception"
	"github.com/erodriguezg/meet/pkg/core/repository"
)

type PersonService interface {
	FindByEmail(email string) (*domain.Person, error)

	FindById(uuid string) (*domain.Person, error)

	FindAll() ([]domain.Person, error)

	FilterPaginated(filters domain.PersonFilter) ([]domain.Person, error)

	Save(person domain.Person) (*domain.Person, error)

	Delete(uuid string) error
}

type domainPersonService struct {
	personRepository repository.PersonRepository
}

func NewDomainPersonService(personRepository repository.PersonRepository) PersonService {
	return &domainPersonService{personRepository}
}

func (port *domainPersonService) FindByEmail(email string) (*domain.Person, error) {
	return port.personRepository.FindByEmail(email)
}

func (port *domainPersonService) FindById(uuid string) (*domain.Person, error) {
	return port.personRepository.FindById(uuid)
}

func (port *domainPersonService) FindAll() ([]domain.Person, error) {
	return port.personRepository.FindAll()
}

func (port *domainPersonService) FilterPaginated(filters domain.PersonFilter) ([]domain.Person, error) {
	return port.personRepository.FilterPaginated(filters)
}

func (port *domainPersonService) Save(person domain.Person) (*domain.Person, error) {

	// validate unique email

	existingPersonWithEmail, err := port.personRepository.FindByEmail(person.Email)
	if err != nil {
		return nil, err
	}
	if existingPersonWithEmail != nil && (person.Id == nil || *person.Id != *existingPersonWithEmail.Id) {
		exception := exception.NewPersonEmailNotAvailableException(&person)
		return nil, exception
	}

	// transaction for save person

	var updatedPerson *domain.Person
	if person.Id == nil {
		person.Active = true
		updatedPerson, err = port.personRepository.Persist(person)
	} else {
		updatedPerson, err = port.personRepository.Update(person)
	}

	return updatedPerson, err
}

func (port *domainPersonService) Delete(uuid string) error {
	return port.personRepository.Delete(uuid)
}
