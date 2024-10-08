package mongodb

import (
	"context"
	"fmt"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ownedResourceCollection = "ownedResources"
)

type ownedResourceMongoDB struct {
	mongoDB *mongo.Database
}

func NewOwnedResourceMongoDB(mongoDB *mongo.Database) repository.OwnedResourceRepository {
	return &ownedResourceMongoDB{mongoDB}
}

func (port *ownedResourceMongoDB) FindByPersonId(personId string) (*domain.OwnedResources, error) {
	personObjectId, _ := primitive.ObjectIDFromHex(personId)
	var ownedResources domain.OwnedResources
	err := port.getCollection().
		FindOne(context.Background(), bson.M{"personId": personObjectId}).
		Decode(&ownedResources)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error OwnedResource FindByPersonId. personId: %s. error: %w", personId, err)
	}
	return &ownedResources, nil
}

func (port *ownedResourceMongoDB) Save(resources domain.OwnedResources) (*domain.OwnedResources, error) {
	filter := bson.M{
		"personId": resources.PersonId,
	}
	update := bson.M{
		"$set": resources,
	}
	opts := options.Update().SetUpsert(true)
	result, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, err
	}
	if resources.Id == nil {
		auxObjectId := result.UpsertedID.(primitive.ObjectID)
		resources.Id = &auxObjectId
	}
	return &resources, nil
}

// private
func (port *ownedResourceMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(ownedResourceCollection)
}
