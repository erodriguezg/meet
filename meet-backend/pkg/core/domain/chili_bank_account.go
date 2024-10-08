package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChiliBankAccount struct {
	Id            *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ModelId       primitive.ObjectID  `json:"modelId" bson:"modelId"`
	Rut           int                 `json:"rut" bson:"rut"`
	HolderName    string              `json:"holderName" bson:"holderName"`
	BankName      string              `json:"bankName" bson:"bankName"`
	AccountType   string              `json:"accountType" bson:"accountType"`
	AccountNumber int                 `json:"accountNumber" bson:"accountNumber"`
}
