package models

type Address struct {
	Id      int    `json:"id"`
	UserId  int    `json:"user_id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	PinCode string `json:"pin_code"`
}
