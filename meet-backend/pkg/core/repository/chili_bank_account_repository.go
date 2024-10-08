package repository

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChiliBankAccountRepository interface {
	FindOneByIdAndModelId(id primitive.ObjectID, modelId primitive.ObjectID) (*domain.ChiliBankAccount, error)
	FindByModelId(modelId primitive.ObjectID) ([]domain.ChiliBankAccount, error)
	Persist(domain.ChiliBankAccount) (*primitive.ObjectID, error)
	Update(domain.ChiliBankAccount) (*primitive.ObjectID, error)
	Delete(id primitive.ObjectID) error
}
