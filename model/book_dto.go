package model

// Book Requests
type CreateBookRequest struct {
	NamaBarang string `json:"nama_barang" validate:"required"`
	Stok       int    `json:"stok" validate:"min=0"`
	Harga      int    `json:"harga" validate:"required,min=0"`
	Keterangan string `json:"keterangan"`
	GambarBuku string `json:"gambar_buku"`
}

type UpdateBookRequest struct {
	NamaBarang string `json:"nama_barang" validate:"required"`
	Stok       int    `json:"stok" validate:"min=0"`
	Terjual    int    `json:"terjual" validate:"min=0"`
	Harga      int    `json:"harga" validate:"required,min=0"`
	Keterangan string `json:"keterangan"`
	GambarBuku string `json:"gambar_buku"`
}
