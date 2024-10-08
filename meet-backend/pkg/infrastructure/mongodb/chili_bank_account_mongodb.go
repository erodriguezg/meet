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
	chiliBankAccountCollection = "chiliBankAccounts"
)

type chiliBankAccountMongoDB struct {
	mongoDB *mongo.Database
}

func NewChiliBankAccountMongoDB(mongoDB *mongo.Database) repository.ChiliBankAccountRepository {
	return &chiliBankAccountMongoDB{mongoDB}
}

func (port *chiliBankAccountMongoDB) Persist(account domain.ChiliBankAccount) (*primitive.ObjectID, error) {
	result, err := port.getCollection().InsertOne(context.Background(), account)
	if err != nil {
		return nil, err
	}
	auxId := result.InsertedID.(primitive.ObjectID)
	return &auxId, nil
}

func (port *chiliBankAccountMongoDB) Update(account domain.ChiliBankAccount) (*primitive.ObjectID, error) {

	filter := bson.M{
		"_id": account.Id,
	}

	update := bson.M{
		"$set": account,
	}

	opts := options.Update().SetUpsert(true)
	_, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, err
	}

	return account.Id, nil
}

func (port *chiliBankAccountMongoDB) Delete(id primitive.ObjectID) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := port.getCollection().DeleteOne(context.Background(), filter)
	return err
}

func (port *chiliBankAccountMongoDB) FindByModelId(modelId primitive.ObjectID) ([]domain.ChiliBankAccount, error) {

	filter := bson.M{
		"modelId": modelId,
	}

	cursor, err := port.getCollection().Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	var accounts []domain.ChiliBankAccount
	err = cursor.All(context.Background(), &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (port *chiliBankAccountMongoDB) FindOneByIdAndModelId(id primitive.ObjectID, modelId primitive.ObjectID) (*domain.ChiliBankAccount, error) {

	filter := bson.M{
		"_id":     id,
		"modelId": modelId,
	}

	var account domain.ChiliBankAccount
	err := port.getCollection().FindOne(context.Background(), filter).Decode(&account)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

// private

func (port *chiliBankAccountMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(chiliBankAccountCollection)
}
