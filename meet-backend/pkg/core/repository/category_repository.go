package repository

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CategoryRepository interface {
	FindAll() ([]domain.Category, error)
	FindById(primitive.ObjectID) (*domain.Category, error)
	FindByParent(parentCategory domain.Category) ([]domain.Category, error)
	Insert(category domain.Category) (*primitive.ObjectID, error)
	Update(category domain.Category) error
	Delete(category domain.Category) error
}
