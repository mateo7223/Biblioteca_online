package models

import "time"

type Purchase struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	BookID      int       `json:"book_id"`
	Quantity    int       `json:"quantity"`
	Total       float64   `json:"total"`
	PurchasedAt time.Time `json:"purchased_at"`
}
