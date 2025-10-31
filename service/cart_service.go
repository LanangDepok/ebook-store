package service

import (
	"fmt"

	"github.com/LanangDepok/ebook-store/entity"
	"github.com/LanangDepok/ebook-store/model"
	"github.com/LanangDepok/ebook-store/repository"
)

type CartService interface {
	AddToCart(userID int, req model.AddToCartRequest) error
	GetCart(userID int) ([]entity.CartItem, int, error)
	UpdateCartItem(cartID int, req model.UpdateCartRequest) error
	RemoveFromCart(cartID int) error
	ClearCart(userID int) error
}

type cartService struct {
	cartRepo repository.CartRepository
	bookRepo repository.BookRepository
}

func NewCartService(cartRepo repository.CartRepository, bookRepo repository.BookRepository) CartService {
	return &cartService{
		cartRepo: cartRepo,
		bookRepo: bookRepo,
	}
}

func (s *cartService) AddToCart(userID int, req model.AddToCartRequest) error {
	// Check if book exists and has stock
	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return fmt.Errorf("book not found")
	}

	if book.Stok < req.Jumlah {
		return fmt.Errorf("insufficient stock")
	}

	// Check if item already in cart
	existingCart, err := s.cartRepo.FindByUserAndBook(userID, req.BookID)
	if err != nil {
		return fmt.Errorf("failed to check cart: %v", err)
	}

	if existingCart != nil {
		// Update existing cart item
		newQuantity := existingCart.Jumlah + req.Jumlah
		if book.Stok < newQuantity {
			return fmt.Errorf("insufficient stock")
		}
		return s.cartRepo.UpdateQuantity(existingCart.ID, newQuantity)
	}

	// Create new cart item
	cart := &entity.Cart{
		UserID: userID,
		BookID: req.BookID,
		Jumlah: req.Jumlah,
		Harga:  book.Harga,
	}

	return s.cartRepo.Create(cart)
}

func (s *cartService) GetCart(userID int) ([]entity.CartItem, int, error) {
	items, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get cart: %v", err)
	}

	total, err := s.cartRepo.GetTotal(userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to calculate total: %v", err)
	}

	return items, total, nil
}

func (s *cartService) UpdateCartItem(cartID int, req model.UpdateCartRequest) error {
	// Note: In a more robust implementation, we should validate the cart belongs to the user
	// and check book stock availability
	return s.cartRepo.UpdateQuantity(cartID, req.Jumlah)
}

func (s *cartService) RemoveFromCart(cartID int) error {
	return s.cartRepo.Delete(cartID)
}

func (s *cartService) ClearCart(userID int) error {
	return s.cartRepo.DeleteByUserID(userID)
}
