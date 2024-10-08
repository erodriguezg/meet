package dto

type PackPaymentMethodDTO struct {
	ChiliBankReceiptMethodEnabled bool     `json:"chiliBankReceiptMethodEnabled" validate:"required"`
	ChiliBankReceiptAccountId     *string  `json:"chiliBankReceiptAccountId,omitempty" validate:"max=255"`
	ChiliBankReceiptCLPPrice      *int     `json:"chiliBankReceiptCLPPrice,omitempty" validate:"required,min=1"`
	PaypalReceiptMethodEnabled    bool     `json:"paypalReceiptMethodEnabled" validate:"required"`
	PaypalReceiptRecipientEmail   *string  `json:"paypalReceiptRecipientEmail,omitempty" validate:"email"`
	PaypalReceiptUSDPrice         *float64 `json:"paypalReceiptUSDPrice,omitempty" validate:"required,min=1"`
	PaypalOnlineMethodEnabled     bool     `json:"paypalOnlineMethodEnabled" validate:"required"`
	PaypalOnlineRecipientEmail    *string  `json:"paypalOnlineRecipientEmail,omitempty" validate:"email"`
	PaypalOnlineUSDPrice          *float64 `json:"paypalOnlineUSDPrice,omitempty" validate:"required,min=1"`
}
