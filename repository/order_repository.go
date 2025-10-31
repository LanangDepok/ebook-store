package repository

import (
	"database/sql"
	"fmt"

	"github.com/LanangDepok/ebook-store/entity"
)

type OrderRepository interface {
	Create(order *entity.Order) error
	CreateItem(item *entity.OrderItem) error
	FindByUserID(userID int) ([]entity.Order, error)
	FindByID(id int) (*entity.OrderDetail, error)
	UpdateStatus(id int, status string) error
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *entity.Order) error {
	query := `
		INSERT INTO orders (user_id, total_harga, status)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, order.UserID, order.TotalHarga, order.Status).
		Scan(&order.ID, &order.CreatedAt)
}

func (r *orderRepository) CreateItem(item *entity.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, book_id, jumlah, harga)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, item.OrderID, item.BookID, item.Jumlah, item.Harga).
		Scan(&item.ID, &item.CreatedAt)
}

func (r *orderRepository) FindByUserID(userID int) ([]entity.Order, error) {
	query := `
		SELECT id, user_id, total_harga, status, created_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		err := rows.Scan(
			&order.ID, &order.UserID, &order.TotalHarga,
			&order.Status, &order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *orderRepository) FindByID(id int) (*entity.OrderDetail, error) {
	// Get order info
	orderQuery := `
		SELECT o.id, o.user_id, o.total_harga, o.status, o.created_at, u.username
		FROM orders o
		JOIN users u ON o.user_id = u.id
		WHERE o.id = $1
	`
	detail := &entity.OrderDetail{}
	err := r.db.QueryRow(orderQuery, id).Scan(
		&detail.ID, &detail.UserID, &detail.TotalHarga,
		&detail.Status, &detail.CreatedAt, &detail.Username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, err
	}

	// Get order items
	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.book_id, oi.jumlah, oi.harga, oi.created_at
		FROM order_items oi
		WHERE oi.order_id = $1
	`
	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.OrderItem
	for rows.Next() {
		var item entity.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.BookID,
			&item.Jumlah, &item.Harga, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	detail.Items = items

	return detail, nil
}

func (r *orderRepository) UpdateStatus(id int, status string) error {
	query := `
		UPDATE orders
		SET status = $1
		WHERE id = $2
	`
	result, err := r.db.Exec(query, status, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}
