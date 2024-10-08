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
	paymentOrderCollection = "paymentsOrders"
)

type paymentOrderMongoDB struct {
	mongoDB *mongo.Database
}

func NewPaymentOrderMongoDB(mongoDB *mongo.Database) repository.PaymentOrderRepository {
	return &paymentOrderMongoDB{mongoDB}
}

func (port *paymentOrderMongoDB) SavePaymentOrder(paymentOrder *domain.PaymentOrder) (*domain.PaymentOrder, error) {
	filter := bson.M{
		"orderId": paymentOrder.OrderId,
	}
	update := bson.M{
		"$set": *paymentOrder,
	}
	opts := options.Update().SetUpsert(true)
	result, err := port.getCollection().UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		return nil, err
	}
	if paymentOrder.Id == nil {
		auxObjectId := result.UpsertedID.(primitive.ObjectID)
		paymentOrder.Id = &auxObjectId
	}
	return paymentOrder, nil
}

func (port *paymentOrderMongoDB) FindByOrderId(orderId string) (*domain.PaymentOrder, error) {
	filter := bson.M{
		"orderId": orderId,
	}
	var paymentOrder domain.PaymentOrder
	err := port.getCollection().FindOne(context.Background(), filter).Decode(&paymentOrder)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &paymentOrder, nil
}

// private

func (port *paymentOrderMongoDB) getCollection() *mongo.Collection {
	return port.mongoDB.Collection(paymentOrderCollection)
}
