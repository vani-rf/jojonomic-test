package models

type Response struct {
	IsError bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    struct {
		HargaTopup   float64 `json:"harga_topup"`
		HargaBuyback float64 `json:"harga_buyback"`
	} `json:"data,omitempty"`
}
