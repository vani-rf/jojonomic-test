package models

type InputHargaRequest struct {
	AdminId      string  `json:"admin_id"`
	ReffId       string  `json:"reff_id,omitempty"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
}
