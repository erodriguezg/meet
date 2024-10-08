package mongodb

import (
	"context"

	"github.com/erodriguezg/meet/pkg/core/domain"
	"github.com/erodriguezg/meet/pkg/core/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	packPaymentMethodCollection = "packsPaymentMethods"
)

type packPaymentMethodMongoDB struct {
	mongoDB *mongo.Database
}

func NewPackPaymentMethodRepository(mongoDB *mongo.Database) repository.PackPaymentMethodRepository {
	return &packPaymentMethodMongoDB{mongoDB}
}

func (port *packPaymentMethodMongoDB) FindByPackId(packId primitive.ObjectID) (*domain.PackPaymentMethod, error) {
	filter := bson.M{
		"packId": packId,
	}
	var packPaymentMethod domain.PackPaymentMethod
	err := port.getCollection().FindOne(context.Background(), filter).Decode(&packPaymentMethod)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &packPaymentMethod, nil
}

func (port *packPaymentMethodMongoDB) Save(paymentMethod domain.PackPaymentMethod) (domain.PackPaymentMethod, error) {
	var filter primitive.M
	if paymentMethod.Id != nil {
		filter = bson.M{
			"_id": paymentMethod.Id,
		}
	} else {
		filter = bson.M{}
	}

	update := bson.M{
		"$set": paymentMethod,
	}

	opts := options.Update().SetUpsert(true)
	result, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return domain.PackPaymentMethod{}, err
	}

	if result.UpsertedID != nil {
		editedId := result.UpsertedID.(primitive.ObjectID)
		paymentMethod.Id = &editedId
	}

	return paymentMethod, nil
}

// private

func (port *packPaymentMethodMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(packPaymentMethodCollection)
}
