package repository

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PackPaymentMethodRepository interface {
	Save(paymentMethod domain.PackPaymentMethod) (domain.PackPaymentMethod, error)

	FindByPackId(packId primitive.ObjectID) (*domain.PackPaymentMethod, error)
}
