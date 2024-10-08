package repository

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileMetaDataRepository interface {
	FindByHash(hash string) (*domain.FileMetaData, error)

	Save(fileMetaData domain.FileMetaData) (*domain.FileMetaData, error)

	Delete(id primitive.ObjectID) error
}
