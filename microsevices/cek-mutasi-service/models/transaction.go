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

type ListTransaction []*Transaction

func (lt ListTransaction) ToResponseItems() []*TransactionResponseItem {
	list := make([]*TransactionResponseItem, len(lt))
	for i, v := range lt {
		list[i] = &TransactionResponseItem{
			CreatedAt:    int32(v.CreatedAt),
			Type:         v.Type,
			GoldWeight:   v.GoldWeight,
			GoldBalance:  v.GoldBalance,
			HargaTopup:   v.HargaTopup,
			HargaBuyback: v.HargaBuyback,
		}
	}

	return list
}
