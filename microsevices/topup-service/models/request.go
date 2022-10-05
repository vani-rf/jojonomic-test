package models

type TopupRequest struct {
	GoldWeight float32 `json:"gram"`
	Amount     float64 `json:"harga"`
	Norek      string  `json:"norek"`
}
