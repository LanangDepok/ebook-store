package model

// Cart Requests
type AddToCartRequest struct {
	BookID int `json:"book_id" validate:"required"`
	Jumlah int `json:"jumlah" validate:"required,min=1"`
}

type UpdateCartRequest struct {
	Jumlah int `json:"jumlah" validate:"required,min=1"`
}
