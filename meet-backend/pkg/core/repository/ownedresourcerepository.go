package repository

import "github.com/erodriguezg/meet/pkg/core/domain"

type OwnedResourceRepository interface {
	FindByPersonId(personId string) (*domain.OwnedResources, error)
	Save(resources domain.OwnedResources) (*domain.OwnedResources, error)
}
