package repository

type PaymentClientRepository interface {
	GetClientData() (map[string]any, error)

	CreateOrder(value float64, currencyCode string) (string, error)

	CapturePayment(orderID string) (map[string]any, error)
}
