package entity

import "time"

type OrderDetail struct {
	ID         int         `json:"id"`
	UserID     int         `json:"user_id"`
	Username   string      `json:"username"`
	TotalHarga int         `json:"total_harga"`
	Status     string      `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	Items      []OrderItem `json:"items"`
}
