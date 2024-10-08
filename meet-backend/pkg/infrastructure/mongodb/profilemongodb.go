package mongodb

import (
	"context"
	"fmt"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	profileCollection = "profiles"
)

type profileMongoDb struct {
	mongoDB *mongo.Database
}

func NewProfileMongoDB(mongoDB *mongo.Database) repository.ProfileRepository {
	return &profileMongoDb{mongoDB}
}

func (port *profileMongoDb) FindByCode(code int) (*domain.Profile, error) {
	var profileAux domain.Profile
	err := port.getCollection().FindOne(context.Background(), bson.M{"code": code}).Decode(&profileAux)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error decoding profile on FindByCode: %w", err)
	}
	return &profileAux, nil
}

func (port *profileMongoDb) FindAll() ([]domain.Profile, error) {
	cursor, err := port.getCollection().Find(context.Background(), bson.M{})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []domain.Profile{}, nil
		}
		return nil, fmt.Errorf("error on FindAll: %w", err)
	}

	var profilesAux []domain.Profile
	err = cursor.All(context.Background(), &profilesAux)
	if err != nil {
		return nil, fmt.Errorf("error on FindAll at decoding results: %w", err)
	}

	return profilesAux, nil
}

// private

func (port *profileMongoDb) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(profileCollection)
}
