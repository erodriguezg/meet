package domain

import (
	"github.com/erodriguezg/meet/pkg/util/datetime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	Id          *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfileCode int                 `json:"profileCode" bson:"profileCode"`
	Email       string              `json:"email" bson:"email"`
	FirstName   string              `json:"firstName" bson:"firstName"`
	LastName    string              `json:"lastName" bson:"lastName"`
	BirthDay    *datetime.Date      `json:"birthday,omitempty" bson:"birthday,omitempty"`
	Active      bool                `json:"active" bson:"active"`
}
