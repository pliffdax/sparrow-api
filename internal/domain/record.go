package domain

import "time"

type Record struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	CategoryID int64     `json:"category_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}
