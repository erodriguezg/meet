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
	packCollection = "packs"
)

type packMongoDB struct {
	mongoDB *mongo.Database
}

func NewPackMongoDB(mongoDB *mongo.Database) repository.PackRepository {
	return &packMongoDB{mongoDB}
}

func (port *packMongoDB) FindPackById(packId string) (*domain.Pack, error) {
	packObjectId, _ := primitive.ObjectIDFromHex(packId)
	var pack domain.Pack
	err := port.getCollection().
		FindOne(context.Background(), bson.M{"_id": packObjectId}).
		Decode(&pack)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error Pack FindPackById. packId: %s. error: %w", packId, err)
	}
	return &pack, nil
}

func (port *packMongoDB) FindPackByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error) {
	modelObjectId, err := primitive.ObjectIDFromHex(modelId)
	if err != nil {
		return nil, fmt.Errorf("error on FindPackByModelIdAndPackNumber getting objectIdFromHex from modelId: %s. error: %w", modelId, err)
	}
	filter := bson.M{
		"modelId":    modelObjectId,
		"packNumber": packNumber,
	}
	var pack domain.Pack
	err = port.getCollection().
		FindOne(context.Background(), filter).
		Decode(&pack)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error Pack FindPackByModelIdAndPackNumber. modelId: %s, packNumber: %d. error: %w", modelId, packNumber, err)
	}
	return &pack, nil
}

func (port *packMongoDB) FindPackActiveByModelIdAndPackNumber(modelId string, packNumber int) (*domain.Pack, error) {
	modelObjectId, err := primitive.ObjectIDFromHex(modelId)
	if err != nil {
		return nil, fmt.Errorf("error on FindPackActiveByModelIdAndPackNumber getting objectIdFromHex from modelId: %s. error: %w", modelId, err)
	}
	filter := bson.M{
		"modelId":    modelObjectId,
		"packNumber": packNumber,
		"active":     true,
	}
	var pack domain.Pack
	err = port.getCollection().
		FindOne(context.Background(), filter).
		Decode(&pack)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error Pack FindPackActiveByModelIdAndPackNumber. modelId: %s, packNumber: %d. error: %w", modelId, packNumber, err)
	}
	return &pack, nil
}

func (port *packMongoDB) FindPacksActiveByModelId(modelId string) ([]domain.Pack, error) {
	modelObjectId, err := primitive.ObjectIDFromHex(modelId)
	if err != nil {
		return nil, fmt.Errorf("error on FindPacksActiveByModelId getting objectIdFromHex from modelId: %s. error: %w", modelId, err)
	}
	filter := bson.M{
		"modelId": modelObjectId,
		"active":  true,
	}

	cursor, err := port.getCollection().
		Find(context.Background(), filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error Pack FindPacksActiveByModelId. modelId: %s. error: %w", modelId, err)
	}

	var packs []domain.Pack
	err = cursor.All(context.Background(), &packs)
	if err != nil {
		return nil, fmt.Errorf("error on FindPacksActiveByModelId decode cursor all. error: %w", err)
	}
	return packs, nil
}

func (port *packMongoDB) SavePack(pack domain.Pack) (*domain.Pack, error) {
	filter := bson.M{
		"modelId":    pack.ModelId,
		"packNumber": pack.PackNumber,
		"active":     pack.Active,
	}

	update := bson.M{
		"$set": pack,
	}

	opts := options.Update().SetUpsert(true)
	result, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, err
	}
	if pack.Id == nil {
		auxObjectId := result.UpsertedID.(primitive.ObjectID)
		pack.Id = &auxObjectId
	}
	return &pack, nil
}

// private

func (port *packMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(packCollection)
}
