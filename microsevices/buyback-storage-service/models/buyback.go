package models

type BuybackParams struct {
	GoldWeight         float32 `json:"gram"`
	Amount             float64 `json:"harga"`
	Norek              string  `json:"norek"`
	ReffID             string  `json:"reff_id,omitempty"`
	HargaTopup         float64 `json:"harga_topup"`
	HargaBuyback       float64 `json:"harga_buyback"`
	CurrentGoldBalance float32 `json:"current_gold_balance"`
}
