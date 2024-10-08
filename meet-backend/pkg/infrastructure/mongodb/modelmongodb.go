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
	modelCollection = "models"
)

type modelMongoDB struct {
	mongoDB *mongo.Database
}

func NewModelMongoDB(mongoDB *mongo.Database) repository.ModelRepository {
	return &modelMongoDB{mongoDB}
}

func (port *modelMongoDB) SearchModelsCount(filters domain.FilterSearchModel) (int, error) {
	filterBson := port.searchFilterBson(filters)
	count, err := port.getCollection().CountDocuments(context.Background(), filterBson)
	if err != nil {
		return -1, err
	}
	return int(count), nil
}

func (port *modelMongoDB) SearchModels(filters domain.FilterSearchModel, first int, last int) ([]domain.Model, error) {

	filterBson := port.searchFilterBson(filters)

	findOptions := options.Find()
	findOptions.SetSkip(int64(first))
	findOptions.SetLimit(int64(last - first))

	cursor, err := port.getCollection().Find(context.Background(), filterBson, findOptions)
	if err != nil {
		return nil, err
	}

	var models []domain.Model
	err = cursor.All(context.Background(), &models)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (port *modelMongoDB) FindModelByPersonId(personId string) (*domain.Model, error) {
	personObjectId, _ := primitive.ObjectIDFromHex(personId)
	var model domain.Model
	err := port.getCollection().
		FindOne(context.Background(), bson.M{"personId": personObjectId}).
		Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error FindModelByPersonId. personId: %s. error: %w", personId, err)
	}
	return &model, nil
}

func (port *modelMongoDB) FindModelByNickName(nickName string) (*domain.Model, error) {
	var model domain.Model
	err := port.getCollection().
		FindOne(context.Background(), bson.M{"nickName": nickName}).
		Decode(&model)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error FindModelByNickName. nickName: %s. error: %w", nickName, err)
	}
	return &model, nil
}

func (port *modelMongoDB) SaveModel(model domain.Model) (*domain.Model, error) {

	filter := bson.M{
		"personId": model.PersonId,
	}

	update := bson.M{
		"$set": model,
	}

	opts := options.Update().SetUpsert(true)

	result, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, err
	}

	if model.Id == nil {
		auxObjectId := result.UpsertedID.(primitive.ObjectID)
		model.Id = &auxObjectId
	}

	return &model, nil
}

// private

func (port *modelMongoDB) searchFilterBson(filters domain.FilterSearchModel) primitive.M {
	filterBson := bson.M{}

	if filters.NickName != nil {
		filterBson["nickName"] = primitive.Regex{Pattern: *filters.NickName, Options: "i"}
	}

	if filters.CityName != nil {
		filterBson["city"] = primitive.Regex{Pattern: *filters.CityName, Options: "i"}
	}

	if filters.CountryCode != nil {
		filterBson["countryCode"] = *filters.CountryCode
	}

	if filters.ZodiacSignCode != nil {
		filterBson["zodiacSignCode"] = *filters.ZodiacSignCode
	}

	return filterBson
}

func (port *modelMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(modelCollection)
}
