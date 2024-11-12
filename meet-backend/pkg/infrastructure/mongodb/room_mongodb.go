package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	roomCollection = "rooms"
)

type roomMongoDB struct {
	mongoDB *mongo.Database
}

func NewRoomMongoDB(mongoDB *mongo.Database) repository.RoomRepository {
	return &roomMongoDB{mongoDB}
}

// FindAll implements repository.RoomRepository.
func (port *roomMongoDB) FindAll() ([]domain.Room, error) {
	filter := bson.M{}
	return findMany[domain.Room](context.Background(), port.getCollection(), filter)
}

// FindByRoomHash implements repository.RoomRepository.
func (port *roomMongoDB) FindByRoomHash(roomHash string) (*domain.Room, error) {
	filter := bson.M{"roomHash": roomHash}
	return findOne[domain.Room](context.Background(), port.getCollection(), filter)
}

// FindById implements repository.RoomRepository.
func (port *roomMongoDB) FindById(roomId primitive.ObjectID) (*domain.Room, error) {
	filter := bson.M{"_id": roomId}
	return findOne[domain.Room](context.Background(), port.getCollection(), filter)
}

// FindByPersonId implements repository.RoomRepository.
func (port *roomMongoDB) FindByOwnerPersonId(ownerPersonId primitive.ObjectID) ([]domain.Room, error) {
	filter := bson.M{"ownerPersonId": ownerPersonId}
	return findMany[domain.Room](context.Background(), port.getCollection(), filter)
}

// FindsWithoutInteractionSince implements repository.RoomRepository.
func (port *roomMongoDB) FindsWithoutInteractionSince(limitDate time.Time) ([]domain.Room, error) {
	filter := bson.M{"lastInteractionDate": bson.M{"$lt": limitDate}}
	return findMany[domain.Room](context.Background(), port.getCollection(), filter)
}

// Persist implements repository.RoomRepository.
func (port *roomMongoDB) Persist(room domain.Room) (*domain.Room, error) {
	result, err := port.getCollection().InsertOne(context.Background(), room)
	if err != nil {
		return nil, err
	}

	auxId, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert InsertedID to ObjectID")
	}
	room.Id = &auxId
	return &room, nil
}

// Update implements repository.RoomRepository.
func (port *roomMongoDB) Update(room domain.Room) (*domain.Room, error) {
	if room.Id == nil {
		return nil, fmt.Errorf("room ID is nil")
	}

	filter := bson.M{"_id": room.Id}
	update := bson.M{"$set": room}
	_, err := port.getCollection().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

// Delete implements repository.RoomRepository.
func (port *roomMongoDB) Delete(roomId primitive.ObjectID) error {
	filter := bson.M{"_id": roomId}
	_, err := port.getCollection().DeleteOne(context.Background(), filter)
	return err
}

// Deletes implements repository.RoomRepository.
func (port *roomMongoDB) Deletes(roomIdList []primitive.ObjectID) error {
	filter := bson.M{"_id": bson.M{"$in": roomIdList}}
	_, err := port.getCollection().DeleteMany(context.Background(), filter)
	return err
}

// private

func (port *roomMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(roomCollection)
}
