package dto

type PackBuyDetailDto struct {
	ModelNickName   string  `json:"modelNickName"`
	PackTitle       *string `json:"packTitle,omitempty"`
	PackDollarValue float64 `json:"packDollarValue"`
}
