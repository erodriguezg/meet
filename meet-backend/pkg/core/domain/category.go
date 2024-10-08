package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Category struct {
	Id       *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ParentId *primitive.ObjectID `json:"parentId,omitempty" bson:"parentId,omitempty"`
	Name     string              `json:"name" bson:"name"`
}
