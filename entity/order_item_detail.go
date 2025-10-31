package entity

type OrderItemDetail struct {
	BookID     int    `json:"book_id"`
	NamaBarang string `json:"nama_barang"`
	Jumlah     int    `json:"jumlah"`
	Harga      int    `json:"harga"`
	Subtotal   int    `json:"subtotal"`
}
