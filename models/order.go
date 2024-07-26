package models

type Order struct {
	Id         int     `json:"id"`
	TxnId      int     `json:"txn_id"`
	UserId     int     `json:"user_id"`
	ProductId  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	Status     string  `json:"status"`
}
