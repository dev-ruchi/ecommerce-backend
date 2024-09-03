package models

type Review struct {
	Id      int     `json:"id"`
	Rating  float64 `json:"rating"`
	Comment string  `json:"comment"`
}
