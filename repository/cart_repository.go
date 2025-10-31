package repository

import (
	"database/sql"
	"fmt"

	"github.com/LanangDepok/ebook-store/entity"
)

type CartRepository interface {
	Create(cart *entity.Cart) error
	FindByUserID(userID int) ([]entity.CartItem, error)
	FindByUserAndBook(userID, bookID int) (*entity.Cart, error)
	UpdateQuantity(id int, quantity int) error
	Delete(id int) error
	DeleteByUserID(userID int) error
	GetTotal(userID int) (int, error)
}

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(cart *entity.Cart) error {
	query := `
		INSERT INTO carts (user_id, book_id, jumlah, harga)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, cart.UserID, cart.BookID, cart.Jumlah, cart.Harga).
		Scan(&cart.ID, &cart.CreatedAt, &cart.UpdatedAt)
}

func (r *cartRepository) FindByUserID(userID int) ([]entity.CartItem, error) {
	query := `
		SELECT 
			c.id, c.book_id, c.jumlah, c.harga,
			b.nama_barang, b.stok, b.harga as harga_satuan,
			COALESCE(b.gambar_buku, '') as gambar_buku
		FROM carts c
		JOIN books b ON c.book_id = b.id
		WHERE c.user_id = $1
		ORDER BY c.created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.CartItem
	for rows.Next() {
		var item entity.CartItem
		err := rows.Scan(
			&item.ID, &item.BookID, &item.Jumlah, &item.Harga,
			&item.NamaBarang, &item.Stok, &item.HargaSatuan, &item.GambarBuku,
		)
		if err != nil {
			return nil, err
		}
		item.Subtotal = item.Jumlah * item.Harga
		items = append(items, item)
	}
	return items, nil
}

func (r *cartRepository) FindByUserAndBook(userID, bookID int) (*entity.Cart, error) {
	query := `
		SELECT id, user_id, book_id, jumlah, harga, created_at, updated_at
		FROM carts
		WHERE user_id = $1 AND book_id = $2
	`
	cart := &entity.Cart{}
	err := r.db.QueryRow(query, userID, bookID).Scan(
		&cart.ID, &cart.UserID, &cart.BookID, &cart.Jumlah,
		&cart.Harga, &cart.CreatedAt, &cart.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return cart, nil
}

func (r *cartRepository) UpdateQuantity(id int, quantity int) error {
	query := `
		UPDATE carts
		SET jumlah = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	result, err := r.db.Exec(query, quantity, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}
	return nil
}

func (r *cartRepository) Delete(id int) error {
	query := `DELETE FROM carts WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}
	return nil
}

func (r *cartRepository) DeleteByUserID(userID int) error {
	query := `DELETE FROM carts WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *cartRepository) GetTotal(userID int) (int, error) {
	query := `
		SELECT COALESCE(SUM(jumlah * harga), 0)
		FROM carts
		WHERE user_id = $1
	`
	var total int
	err := r.db.QueryRow(query, userID).Scan(&total)
	return total, err
}
