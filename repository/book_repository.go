package repository

import (
	"database/sql"
	"fmt"

	"github.com/LanangDepok/ebook-store/entity"
)

type BookRepository interface {
	Create(book *entity.Book) error
	FindAll() ([]entity.Book, error)
	FindByID(id int) (*entity.Book, error)
	Update(id int, book *entity.Book) error
	Delete(id int) error
	UpdateStock(id int, quantity int) error
	IncrementSold(id int, quantity int) error
}

type bookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *entity.Book) error {
	query := `
		INSERT INTO books (nama_barang, stok, harga, keterangan, gambar_buku)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, terjual, created_at, updated_at
	`
	return r.db.QueryRow(query, book.NamaBarang, book.Stok, book.Harga,
		book.Keterangan, book.GambarBuku).
		Scan(&book.ID, &book.Terjual, &book.CreatedAt, &book.UpdatedAt)
}

func (r *bookRepository) FindAll() ([]entity.Book, error) {
	query := `
		SELECT id, nama_barang, stok, terjual, harga, keterangan,
		       COALESCE(gambar_buku, '') as gambar_buku, created_at, updated_at
		FROM books
		ORDER BY id DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		err := rows.Scan(
			&book.ID, &book.NamaBarang, &book.Stok, &book.Terjual,
			&book.Harga, &book.Keterangan, &book.GambarBuku,
			&book.CreatedAt, &book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *bookRepository) FindByID(id int) (*entity.Book, error) {
	query := `
		SELECT id, nama_barang, stok, terjual, harga, keterangan,
		       COALESCE(gambar_buku, '') as gambar_buku, created_at, updated_at
		FROM books
		WHERE id = $1
	`
	book := &entity.Book{}
	err := r.db.QueryRow(query, id).Scan(
		&book.ID, &book.NamaBarang, &book.Stok, &book.Terjual,
		&book.Harga, &book.Keterangan, &book.GambarBuku,
		&book.CreatedAt, &book.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, err
	}
	return book, nil
}

func (r *bookRepository) Update(id int, book *entity.Book) error {
	query := `
		UPDATE books
		SET nama_barang = $1, stok = $2, terjual = $3, harga = $4,
		    keterangan = $5, gambar_buku = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7
		RETURNING updated_at
	`
	result := r.db.QueryRow(query, book.NamaBarang, book.Stok, book.Terjual,
		book.Harga, book.Keterangan, book.GambarBuku, id)

	err := result.Scan(&book.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("book not found")
		}
		return err
	}
	return nil
}

func (r *bookRepository) Delete(id int) error {
	query := `DELETE FROM books WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("book not found")
	}
	return nil
}

func (r *bookRepository) UpdateStock(id int, quantity int) error {
	query := `
		UPDATE books
		SET stok = stok - $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2 AND stok >= $1
	`
	result, err := r.db.Exec(query, quantity, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("insufficient stock or book not found")
	}
	return nil
}

func (r *bookRepository) IncrementSold(id int, quantity int) error {
	query := `
		UPDATE books
		SET terjual = terjual + $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(query, quantity, id)
	return err
}
