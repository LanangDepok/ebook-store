package service

import (
	"fmt"

	"github.com/LanangDepok/ebook-store/entity"
	"github.com/LanangDepok/ebook-store/model"
	"github.com/LanangDepok/ebook-store/repository"
)

type BookService interface {
	CreateBook(req model.CreateBookRequest) (*entity.Book, error)
	GetAllBooks() ([]entity.Book, error)
	GetBookByID(id int) (*entity.Book, error)
	UpdateBook(id int, req model.UpdateBookRequest) (*entity.Book, error)
	DeleteBook(id int) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) CreateBook(req model.CreateBookRequest) (*entity.Book, error) {
	book := &entity.Book{
		NamaBarang: req.NamaBarang,
		Stok:       req.Stok,
		Harga:      req.Harga,
		Keterangan: req.Keterangan,
		GambarBuku: req.GambarBuku,
		Terjual:    0,
	}

	err := s.repo.Create(book)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %v", err)
	}

	return book, nil
}

func (s *bookService) GetAllBooks() ([]entity.Book, error) {
	books, err := s.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %v", err)
	}
	return books, nil
}

func (s *bookService) GetBookByID(id int) (*entity.Book, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("book not found: %v", err)
	}
	return book, nil
}

func (s *bookService) UpdateBook(id int, req model.UpdateBookRequest) (*entity.Book, error) {
	// Check if book exists
	existingBook, err := s.repo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("book not found: %v", err)
	}

	// Update book fields
	existingBook.NamaBarang = req.NamaBarang
	existingBook.Stok = req.Stok
	existingBook.Terjual = req.Terjual
	existingBook.Harga = req.Harga
	existingBook.Keterangan = req.Keterangan
	if req.GambarBuku != "" {
		existingBook.GambarBuku = req.GambarBuku
	}

	err = s.repo.Update(id, existingBook)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %v", err)
	}

	return existingBook, nil
}

func (s *bookService) DeleteBook(id int) error {
	// Check if book exists
	_, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("book not found: %v", err)
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete book: %v", err)
	}

	return nil
}
