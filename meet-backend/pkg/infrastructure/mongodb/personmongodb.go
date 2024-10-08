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
	personCollection = "persons"
)

type personMongoDB struct {
	mongoDB *mongo.Database
}

func NewPersonMongoDB(mongoDB *mongo.Database) repository.PersonRepository {
	return &personMongoDB{mongoDB}
}

func (*personMongoDB) Delete(uuid string) error {
	panic("unimplemented")
}

func (*personMongoDB) FilterPaginated(filters domain.PersonFilter) ([]domain.Person, error) {
	panic("unimplemented")
}

func (*personMongoDB) FindAll() ([]domain.Person, error) {
	panic("unimplemented")
}

func (port *personMongoDB) FindByEmail(email string) (*domain.Person, error) {
	var person domain.Person
	err := port.getCollection().
		FindOne(context.Background(), bson.M{"email": email}).
		Decode(&person)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error FindByEmail. email: %s, error: %w", email, err)
	}

	return &person, err
}

func (port *personMongoDB) FindById(uuid string) (*domain.Person, error) {
	personObjectId, _ := primitive.ObjectIDFromHex(uuid)
	var person domain.Person
	err := port.getCollection().
		FindOne(context.Background(), bson.M{"_id": personObjectId}).
		Decode(&person)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error FindById. personId: %s. error: %w", uuid, err)
	}
	return &person, nil
}

func (port *personMongoDB) Persist(person domain.Person) (*domain.Person, error) {
	return port.upsert(person)
}

func (port *personMongoDB) Update(person domain.Person) (*domain.Person, error) {
	return port.upsert(person)
}

// private

func (port *personMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(personCollection)
}

func (port *personMongoDB) upsert(person domain.Person) (*domain.Person, error) {
	filter := bson.M{"email": person.Email}
	update := bson.M{"$set": person}
	opts := options.Update().SetUpsert(true)

	result, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, err
	}

	if person.Id == nil {
		upsertedId := result.UpsertedID.(primitive.ObjectID)
		person.Id = &upsertedId
	}
	return &person, nil

}
