package service

import (
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OwnedResourceService interface {
	AddPackToPerson(personId string, packId string) error

	PersonHasPack(personId string, packId string) (bool, error)
}

type domainOwnedResourceService struct {
	repository repository.OwnedResourceRepository
}

func NewDomainOwnedResourceService(repository repository.OwnedResourceRepository) OwnedResourceService {
	return &domainOwnedResourceService{repository}
}

func (port *domainOwnedResourceService) AddPackToPerson(personId string, packId string) error {
	searchResources, err := port.repository.FindByPersonId(personId)
	if err != nil {
		return fmt.Errorf("error at FindByPersonId in AddPackToPerson. personId: %s, error: %w", personId, err)
	}

	actualDate := time.Now()

	personObjectId, err := primitive.ObjectIDFromHex(personId)
	if err != nil {
		return fmt.Errorf("error creating ObjectId in AddPackToPerson. personId: %s, error: %w", personId, err)
	}

	packObjectId, err := primitive.ObjectIDFromHex(packId)
	if err != nil {
		return fmt.Errorf("error creating ObjectId in AddPackToPerson. packId: %s, error: %w", packId, err)
	}

	var resources domain.OwnedResources
	if searchResources == nil {

		resources = domain.OwnedResources{
			PersonId:       personObjectId,
			OwnedPacksId:   []primitive.ObjectID{packObjectId},
			LastUpdateDate: actualDate,
		}
	} else {
		exist := port.searchPackInResources(searchResources, packId)
		if exist {
			return nil
		}
		resources = *searchResources
		resources.OwnedPacksId = append(resources.OwnedPacksId, packObjectId)
		resources.LastUpdateDate = actualDate
	}

	_, err = port.repository.Save(resources)
	if err != nil {
		return fmt.Errorf("error on Save Person Resources Repository at AddPackToPerson. personId: %s, packId: %s, error: %w", personId, packId, err)
	}

	return nil
}

func (port *domainOwnedResourceService) PersonHasPack(personId string, packId string) (bool, error) {
	resources, err := port.repository.FindByPersonId(personId)
	if err != nil {
		return false, fmt.Errorf("error at FindByPersonId in PersonHasPack. personId: %s, error: %w", personId, err)
	}
	if resources == nil {
		return false, nil
	}
	return port.searchPackInResources(resources, packId), nil
}

// private
func (port *domainOwnedResourceService) searchPackInResources(resources *domain.OwnedResources, packId string) bool {
	for _, ownedPackId := range resources.OwnedPacksId {
		if ownedPackId.Hex() == packId {
			return true
		}
	}
	return false
}
