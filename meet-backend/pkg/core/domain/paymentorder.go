package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentOrder struct {
	Id                 *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OrderId            string              `json:"orderId" bson:"orderId"`
	PersonId           primitive.ObjectID  `json:"personId" bson:"personId"`
	PackId             primitive.ObjectID  `json:"packId" bson:"packId"`
	ModelId            primitive.ObjectID  `json:"modelId" bson:"modelId"`
	PaymentDollarValue float64             `json:"paymentDollarValue" bson:"paymentDollarValue"`
	CreatedAt          time.Time           `json:"createdAt" bson:"createdAt"`
	CapturedAt         *time.Time          `json:"capturedAt,omitempty" bson:"capturedAt,omitempty"`
	ModelPaidAt        *time.Time          `json:"modelPaidAt,omitempty" bson:"modelPaidAt,omitempty"`
	PaymentDetails     map[string]any      `json:"paymentDetails,omitempty" bson:"paymentDetails,omitempty"`
}
