package models

import "time"

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Price     float64   `json:"price"`
	Format    string    `json:"format"`
	Available bool      `json:"available"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}
