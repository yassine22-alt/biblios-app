package model

import "time"

type Order struct {
	ID         int         `json:"id"`
	CustomerId int         `json:"customer"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	CreatedAt  time.Time   `json:"created_at"`
	Status     string      `json:"status"`
}

type OrderInput struct {
	CustomerId int         `json:"customer"`
	Items      []OrderItem `json:"items"`
}
