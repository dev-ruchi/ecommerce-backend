package models

import (
	"github.com/lib/pq"
)

type Product struct {
	Id          int            `json:"id"`
	Title       string         `json:"title"`
	Price       float64        `json:"price"`
	Description string         `json:"description"`
	Rating      float64        `json:"rating"`
	Images      pq.StringArray `json:"images"`
}
