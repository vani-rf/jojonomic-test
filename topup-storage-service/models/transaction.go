package models

type Transaction struct {
	ReffID       string  `json:"reff_id"`
	Norek        string  `json:"norek"`
	Type         string  `json:"type"`
	GoldWeight   float32 `json:"gold_weight"`
	HargaTopup   float64 `json:"harga_topup"`
	HargaBuyback float64 `json:"harga_buyback"`
	GoldBalance  float32 `json:"gold_balance"`
	CreatedAt    int     `json:"created_at"`
}

func (Transaction) TableName() string {
	return "tbl_transaksi"
}
