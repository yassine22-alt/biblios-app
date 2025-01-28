package model

import "time"

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	AuthorID    int       `json:"author_id"`
	Genres      []string  `json:"genres"`
	PublishedAt time.Time `json:"published_at"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
}

type BookInput struct {
	Title    string   `json:"title"`
	AuthorID int      `json:"author_id"`
	Genres   []string `json:"genres"`
	Price    float64  `json:"price"`
	Stock    int      `json:"stock"`
}

type BookSale struct {
	BookID   int
	Quantity int
}
