package dto

import (
	"fmt"

	"github.com/erodriguezg/meet/pkg/core/domain"
	objectidutils "github.com/erodriguezg/meet/pkg/util/object_id_utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryDTO struct {
	Id       *string        `json:"id,omitempty"`
	ParentId *string        `json:"parentId,omitempty"`
	Name     string         `json:"name"`
	Children []*CategoryDTO `json:"children"`
}

func (port *CategoryDTO) ToDomain() (domain.Category, error) {
	var catDomain domain.Category
	var categoryId *primitive.ObjectID
	var categoryParentId *primitive.ObjectID
	var err error

	categoryId, err = objectidutils.ObjectIDFromHexOrNil(port.Id)
	if err != nil {
		return catDomain, err
	}

	categoryParentId, err = objectidutils.ObjectIDFromHexOrNil(port.ParentId)
	if err != nil {
		return catDomain, err
	}

	catDomain.Id = categoryId
	catDomain.Name = port.Name
	catDomain.ParentId = categoryParentId

	return catDomain, nil
}

func (port *CategoryDTO) FromDomain(catDomain *domain.Category) error {

	if catDomain == nil {
		return fmt.Errorf("category domain is nil")
	}

	port.Id = objectidutils.HexFromObjectIDOrNil(catDomain.Id)
	port.Name = catDomain.Name
	port.ParentId = objectidutils.HexFromObjectIDOrNil(catDomain.ParentId)

	return nil
}
