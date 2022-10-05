package models

type TransactionResponseItem struct {
	CreatedAt    int32   `json:"date"`
	Type         string  `json:"type"`
	GoldWeight   float32 `json:"gram"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
	GoldBalance  float32 `json:"saldo"`
}

type TransactionResponse struct {
	IsError bool                       `json:"error"`
	Data    []*TransactionResponseItem `json:"data,omitempty"`
	Message string                     `json:"message,omitempty"`
}
