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
	fileMetaDataCollection = "filesMetaData"
)

type fileMetaDataMongoDB struct {
	mongoDB *mongo.Database
}

func NewFileMetaDataMongoDB(mongoDB *mongo.Database) repository.FileMetaDataRepository {
	return &fileMetaDataMongoDB{mongoDB}
}

func (port *fileMetaDataMongoDB) Delete(id primitive.ObjectID) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := port.getCollection().DeleteOne(context.Background(), filter)
	return err
}

func (port *fileMetaDataMongoDB) FindByHash(hash string) (*domain.FileMetaData, error) {
	var fileMetaData domain.FileMetaData
	err := port.getCollection().
		FindOne(context.Background(), bson.M{"hash": hash}).
		Decode(&fileMetaData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error at findOne mongodb. error: %w", err)
	}
	return &fileMetaData, nil
}

func (port *fileMetaDataMongoDB) Save(fileMetaData domain.FileMetaData) (*domain.FileMetaData, error) {
	filter := bson.M{
		"hash": fileMetaData.Hash,
	}

	update := bson.M{
		"$set": fileMetaData,
	}

	opts := options.Update().SetUpsert(true)

	result, err := port.getCollection().
		UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, fmt.Errorf("error at updateOne mongoDb. error: %w", err)
	}

	if fileMetaData.Id == nil {
		auxObjectId := result.UpsertedID.(primitive.ObjectID)
		fileMetaData.Id = &auxObjectId
	}

	return &fileMetaData, nil
}

// private

func (port *fileMetaDataMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(fileMetaDataCollection)
}
