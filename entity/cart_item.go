package entity

type CartItem struct {
	ID          int    `json:"id"`
	BookID      int    `json:"book_id"`
	NamaBarang  string `json:"nama_barang"`
	Jumlah      int    `json:"jumlah"`
	Harga       int    `json:"harga"`
	Stok        int    `json:"stok"`
	HargaSatuan int    `json:"harga_satuan"`
	Subtotal    int    `json:"subtotal"`
	GambarBuku  string `json:"gambar_buku"`
}
