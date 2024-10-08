package exception

import "github.com/erodriguezg/meet/pkg/core/domain"

func NewPersonEmailNotAvailableException(person *domain.Person) error {
	return newBusinessException("person-email-not-available", "the email is not available", map[string]string{"email": person.Email})
}

func NewPersonIsNotActiveException(person *domain.Person) error {
	return newBusinessException("person-is-not-active", "the person is not active", map[string]string{"email": person.Email})
}
