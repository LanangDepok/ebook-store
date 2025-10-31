package service

import (
	"database/sql"
	"fmt"

	"github.com/LanangDepok/ebook-store/entity"
	"github.com/LanangDepok/ebook-store/repository"
)

type OrderService interface {
	CreateOrder(userID int) (*entity.Order, error)
	GetUserOrders(userID int) ([]entity.Order, error)
	GetOrderDetail(orderID, userID int) (*entity.OrderDetail, error)
	UpdateOrderStatus(orderID int, status string) error
}

type orderService struct {
	orderRepo repository.OrderRepository
	cartRepo  repository.CartRepository
	bookRepo  repository.BookRepository
	db        *sql.DB
}

func NewOrderService(orderRepo repository.OrderRepository, cartRepo repository.CartRepository, bookRepo repository.BookRepository, db *sql.DB) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
		bookRepo:  bookRepo,
		db:        db,
	}
}

func (s *orderService) CreateOrder(userID int) (*entity.Order, error) {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	defer tx.Rollback()

	// Get cart items
	cartItems, err := s.cartRepo.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %v", err)
	}

	if len(cartItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Calculate total and validate stock
	totalHarga := 0
	for _, item := range cartItems {
		book, err := s.bookRepo.FindByID(item.BookID)
		if err != nil {
			return nil, fmt.Errorf("book not found: %v", err)
		}

		if book.Stok < item.Jumlah {
			return nil, fmt.Errorf("insufficient stock for %s", book.NamaBarang)
		}

		totalHarga += item.Subtotal
	}

	// Create order
	order := &entity.Order{
		UserID:     userID,
		TotalHarga: totalHarga,
		Status:     "pending",
	}

	err = s.orderRepo.Create(order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %v", err)
	}

	// Create order items and update stock
	for _, item := range cartItems {
		orderItem := &entity.OrderItem{
			OrderID: order.ID,
			BookID:  item.BookID,
			Jumlah:  item.Jumlah,
			Harga:   item.Harga,
		}

		err = s.orderRepo.CreateItem(orderItem)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %v", err)
		}

		// Update book stock
		err = s.bookRepo.UpdateStock(item.BookID, item.Jumlah)
		if err != nil {
			return nil, fmt.Errorf("failed to update stock: %v", err)
		}

		// Increment sold count
		err = s.bookRepo.IncrementSold(item.BookID, item.Jumlah)
		if err != nil {
			return nil, fmt.Errorf("failed to update sold count: %v", err)
		}
	}

	// Clear cart
	err = s.cartRepo.DeleteByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to clear cart: %v", err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %v", err)
	}

	return order, nil
}

func (s *orderService) GetUserOrders(userID int) ([]entity.Order, error) {
	return s.orderRepo.FindByUserID(userID)
}

func (s *orderService) GetOrderDetail(orderID, userID int) (*entity.OrderDetail, error) {
	detail, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, err
	}

	// Verify order belongs to user
	if detail.UserID != userID {
		return nil, fmt.Errorf("unauthorized access to order")
	}

	return detail, nil
}

func (s *orderService) UpdateOrderStatus(orderID int, status string) error {
	return s.orderRepo.UpdateStatus(orderID, status)
}
