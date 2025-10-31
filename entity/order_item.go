package entity

import "time"

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"order_id"`
	BookID    int       `json:"book_id"`
	Jumlah    int       `json:"jumlah"`
	Harga     int       `json:"harga"`
	CreatedAt time.Time `json:"created_at"`
}
