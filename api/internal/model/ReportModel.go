package model

import "time"

type ReportModel struct {
	TotalRevenue    float64   `json:"total_revenue"`
	TotalOrders     int       `json:"total_orders"`
	TotalBooksSold  int       `json:"total_books_sold"`
	TopSellingBooks []Book    `json:"top_selling_books"`
	GeneratedAt     time.Time `json:"generated_at"`
}
