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
	categoryCollection = "categories"
)

type categoryMongoDB struct {
	mongoDB *mongo.Database
}

func NewCategoryMongoDB(mongoDB *mongo.Database) repository.CategoryRepository {
	return &categoryMongoDB{mongoDB}
}

// FindAll implements repository.CategoryRepository.
func (port *categoryMongoDB) FindAll() ([]domain.Category, error) {
	filter := bson.M{}
	findOptions := options.Find()
	findOptions.SetSort(
		bson.M{
			"name": 1,
		},
	)
	return port.findMany(filter, findOptions)
}

// FindByParent implements repository.CategoryRepository.
func (port *categoryMongoDB) FindByParent(parentCategory domain.Category) ([]domain.Category, error) {

	var categories []domain.Category
	if parentCategory.Id == nil {
		return categories, nil
	}

	filter := bson.M{
		"parentId": parentCategory.Id,
	}
	findOptions := options.Find()
	findOptions.SetSort(
		bson.M{
			"name": 1,
		},
	)
	return port.findMany(filter, findOptions)
}

// FindById implements repository.CategoryRepository.
func (port *categoryMongoDB) FindById(categoryId primitive.ObjectID) (*domain.Category, error) {
	filter := bson.M{
		"_id": categoryId,
	}
	return port.findOne(filter)
}

// Delete implements repository.CategoryRepository.
func (port *categoryMongoDB) Delete(category domain.Category) error {
	if category.Id != nil {
		id := category.Id
		filter := bson.M{
			"_id": *id,
		}
		_, err := port.getCollection().DeleteOne(context.Background(), filter)
		if err != nil {
			return err
		}
	}
	return nil
}

// Save implements repository.CategoryRepository.
func (port *categoryMongoDB) Insert(category domain.Category) (*primitive.ObjectID, error) {
	result, err := port.getCollection().InsertOne(context.Background(), category)
	if err != nil {
		return nil, err
	}
	auxId := result.InsertedID.(primitive.ObjectID)
	return &auxId, nil
}

// Update implements repository.CategoryRepository.
func (port *categoryMongoDB) Update(category domain.Category) error {

	filter := bson.M{
		"_id": category.Id,
	}

	update := bson.M{
		"$set": category,
	}

	opts := options.Update()
	_, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

// privates

func (port *categoryMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(categoryCollection)
}

func (port *categoryMongoDB) findMany(filter bson.M, findOptions *options.FindOptions) ([]domain.Category, error) {
	var categories []domain.Category
	cursor, err := port.getCollection().Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &categories)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (port *categoryMongoDB) findOne(filter bson.M) (*domain.Category, error) {
	var category domain.Category
	err := port.getCollection().FindOne(context.Background(), filter).Decode(&category)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}
