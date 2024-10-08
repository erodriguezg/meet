package repository

import (
	"github.com/erodriguezg/meet/pkg/core/domain"
)

type PaymentOrderRepository interface {
	SavePaymentOrder(paymentOrder *domain.PaymentOrder) (*domain.PaymentOrder, error)
	FindByOrderId(orderId string) (*domain.PaymentOrder, error)
}
