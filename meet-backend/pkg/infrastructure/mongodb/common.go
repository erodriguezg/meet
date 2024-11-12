package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func findMany[T any](ctx context.Context, collection *mongo.Collection, filter bson.M) ([]T, error) {
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []T
	err = cursor.All(ctx, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func findOne[T any](ctx context.Context, collection *mongo.Collection, filter bson.M) (*T, error) {
	var result T
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}
