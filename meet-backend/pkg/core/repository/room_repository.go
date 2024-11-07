package repository

import (
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomRepository interface {
	FindAll() ([]domain.Room, error)

	FindById(roomId primitive.ObjectID) (*domain.Room, error)

	FindByRoomHash(roomHash string) (*domain.Room, error)

	FindByOwnerPersonId(ownerPersonId primitive.ObjectID) ([]domain.Room, error)

	FindsWithoutInteractionSince(limitDate time.Time) ([]domain.Room, error)

	Persist(room domain.Room) (*domain.Room, error)

	Update(room domain.Room) (*domain.Room, error)

	Delete(roomId primitive.ObjectID) error

	Deletes(roomIdList []primitive.ObjectID) error
}
