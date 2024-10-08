package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OwnedResources struct {
	Id             *primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	PersonId       primitive.ObjectID   `json:"personId" bson:"personId"`
	OwnedPacksId   []primitive.ObjectID `json:"ownedPacksId" bson:"ownedPacksId"`
	LastUpdateDate time.Time            `json:"lastUpdateDate" bson:"lastUpdateDate"`
}
