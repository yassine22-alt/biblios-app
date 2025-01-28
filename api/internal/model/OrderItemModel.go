package model

type OrderItem struct {
	BookID   int `json:"book_id"`
	Quantity int `json:"quantity"`
}
