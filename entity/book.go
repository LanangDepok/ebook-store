package entity

import "time"

type Book struct {
	ID         int       `json:"id"`
	NamaBarang string    `json:"nama_barang"`
	Stok       int       `json:"stok"`
	Terjual    int       `json:"terjual"`
	Harga      int       `json:"harga"`
	Keterangan string    `json:"keterangan"`
	GambarBuku string    `json:"gambar_buku"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
