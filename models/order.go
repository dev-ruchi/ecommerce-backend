package models

type Order struct {
	OrderId    int     `json:"order_id"`
	UderId     int     `json:"customer_id"`
	ProductId  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
}
