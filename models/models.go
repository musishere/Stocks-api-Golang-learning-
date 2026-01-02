package models

import "time"

type Stock struct {
	stockId   int64     `json:"stockId"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Company   string    `json:"company"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
