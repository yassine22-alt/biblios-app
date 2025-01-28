package model

import "time"

type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Address   Address   `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type CustomerInput struct {
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Address Address `json:"address"`
}
