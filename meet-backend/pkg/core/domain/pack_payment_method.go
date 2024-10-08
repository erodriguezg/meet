package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackPaymentMethod struct {
	Id                            *primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PackId                        primitive.ObjectID  `json:"packId" bson:"packId"`
	ChiliBankReceiptMethodEnabled bool                `json:"chiliBankReceiptMethodEnabled" bson:"chiliBankReceiptMethodEnabled"`
	ChiliBankReceiptAccountId     *primitive.ObjectID `json:"chiliBankReceiptAccountId,omitempty" bson:"chiliBankReceiptAccountId,omitempty"`
	ChiliBankReceiptCLPPrice      *int                `json:"chiliBankReceiptCLPPrice,omitempty" bson:"chiliBankReceiptCLPPrice,omitempty"`
	PaypalReceiptMethodEnabled    bool                `json:"paypalReceiptMethodEnabled" bson:"paypalReceiptMethodEnabled"`
	PaypalReceiptRecipientEmail   *string             `json:"paypalReceiptRecipientEmail,omitempty" bson:"paypalReceiptRecipientEmail,omitempty"`
	PaypalReceiptUSDPrice         *float64            `json:"paypalReceiptUSDPrice,omitempty" bson:"paypalReceiptUSDPrice,omitempty"`
	PaypalOnlineMethodEnabled     bool                `json:"paypalOnlineMethodEnabled" bson:"paypalOnlineMethodEnabled"`
	PaypalOnlineRecipientEmail    *string             `json:"paypalOnlineRecipientEmail,omitempty" bson:"paypalOnlineRecipientEmail,omitempty"`
	PaypalOnlineUSDPrice          *float64            `json:"paypalOnlineUSDPrice,omitempty" bson:"paypalOnlineUSDPrice,omitempty"`
	UpdateDate                    time.Time           `json:"updateDate" bson:"updateDate"`
}
