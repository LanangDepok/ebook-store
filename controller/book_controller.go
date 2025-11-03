package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/LanangDepok/ebook-store/model"
	"github.com/LanangDepok/ebook-store/service"
)

type BookController struct {
	bookService   service.BookService
	uploadService service.UploadService
}

func NewBookController(bookService service.BookService, uploadService service.UploadService) *BookController {
	return &BookController{
		bookService:   bookService,
		uploadService: uploadService,
	}
}

func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form for file upload
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Get form values
	namaBarang := r.FormValue("nama_barang")
	stok := r.FormValue("stok")
	harga := r.FormValue("harga")
	keterangan := r.FormValue("keterangan")

	if namaBarang == "" || harga == "" {
		respondError(w, http.StatusBadRequest, "nama_barang and harga are required")
		return
	}

	// Convert string to int
	stokInt := 0
	if stok != "" {
		fmt.Sscanf(stok, "%d", &stokInt)
	}

	hargaInt := 0
	fmt.Sscanf(harga, "%d", &hargaInt)

	if hargaInt <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid price")
		return
	}

	// Handle image upload
	var gambarBuku string
	file, header, err := r.FormFile("gambar_buku")
	if err == nil {
		defer file.Close()
		gambarBuku, err = c.uploadService.UploadImage(file, header)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	req := model.CreateBookRequest{
		NamaBarang: namaBarang,
		Stok:       stokInt,
		Harga:      hargaInt,
		Keterangan: keterangan,
		GambarBuku: gambarBuku,
	}

	book, err := c.bookService.CreateBook(req)
	if err != nil {
		// Delete uploaded image if book creation fails
		if gambarBuku != "" {
			c.uploadService.DeleteImage(gambarBuku)
		}
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Add image URL to response
	if book.GambarBuku != "" {
		book.GambarBuku = c.uploadService.GetImageURL(book.GambarBuku)
	}

	respondSuccess(w, http.StatusCreated, "Book created successfully", book)
}

func (c *BookController) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := c.bookService.GetAllBooks()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Add image URLs to response
	for i := range books {
		if books[i].GambarBuku != "" {
			books[i].GambarBuku = c.uploadService.GetImageURL(books[i].GambarBuku)
		}
	}

	respondSuccess(w, http.StatusOK, "Books retrieved successfully", books)
}

func (c *BookController) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	book, err := c.bookService.GetBookByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	// Add image URL to response
	if book.GambarBuku != "" {
		book.GambarBuku = c.uploadService.GetImageURL(book.GambarBuku)
	}

	respondSuccess(w, http.StatusOK, "Book retrieved successfully", book)
}

func (c *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	// Get existing book to retrieve old image
	existingBook, err := c.bookService.GetBookByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Get form values
	namaBarang := r.FormValue("nama_barang")
	stok := r.FormValue("stok")
	terjual := r.FormValue("terjual")
	harga := r.FormValue("harga")
	keterangan := r.FormValue("keterangan")

	if namaBarang == "" || harga == "" {
		respondError(w, http.StatusBadRequest, "nama_barang and harga are required")
		return
	}

	// Convert string to int
	stokInt := 0
	if stok != "" {
		fmt.Sscanf(stok, "%d", &stokInt)
	}

	terjualInt := 0
	if terjual != "" {
		fmt.Sscanf(terjual, "%d", &terjualInt)
	}

	hargaInt := 0
	fmt.Sscanf(harga, "%d", &hargaInt)

	if hargaInt <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid price")
		return
	}

	// Handle image upload
	gambarBuku := existingBook.GambarBuku
	file, header, err := r.FormFile("gambar_buku")
	if err == nil {
		defer file.Close()

		// Upload new image
		newImage, err := c.uploadService.UploadImage(file, header)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Delete old image if exists
		if existingBook.GambarBuku != "" {
			c.uploadService.DeleteImage(existingBook.GambarBuku)
		}

		gambarBuku = newImage
	}

	req := model.UpdateBookRequest{
		NamaBarang: namaBarang,
		Stok:       stokInt,
		Terjual:    terjualInt,
		Harga:      hargaInt,
		Keterangan: keterangan,
		GambarBuku: gambarBuku,
	}

	book, err := c.bookService.UpdateBook(id, req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Add image URL to response
	if book.GambarBuku != "" {
		book.GambarBuku = c.uploadService.GetImageURL(book.GambarBuku)
	}

	respondSuccess(w, http.StatusOK, "Book updated successfully", book)
}

func (c *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Book ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid book ID")
		return
	}

	// Get book to retrieve image before deletion
	book, err := c.bookService.GetBookByID(id)
	if err != nil {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	// Delete book
	err = c.bookService.DeleteBook(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Delete associated image
	if book.GambarBuku != "" {
		c.uploadService.DeleteImage(book.GambarBuku)
	}

	respondSuccess(w, http.StatusOK, "Book deleted successfully", nil)
}
