package entity

import "time"

type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	TotalHarga int       `json:"total_harga"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
