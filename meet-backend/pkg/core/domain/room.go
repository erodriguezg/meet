package domain

import (
	"github.com/erodriguezg/meet/pkg/util/datetime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	Id                  *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	OwnerPersonId       primitive.ObjectID  `json:"ownerPersonId" bson:"ownerPersonId"`
	CreationDate        datetime.Date       `json:"creationDate" bson:"creationDate"`
	AnonymousAccess     bool                `json:"anonymousAccess" bson:"anonymousAccess"`
	LastInteractionDate datetime.Date       `json:"lastInteractionDate" bson:"lastInteractionDate"`
}
