package controller

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/LanangDepok/ebook-store/service"
)

type UploadController struct {
	uploadService service.UploadService
	uploadDir     string
}

func NewUploadController(uploadService service.UploadService, uploadDir string) *UploadController {
	return &UploadController{
		uploadService: uploadService,
		uploadDir:     uploadDir,
	}
}

// UploadImage handles image upload
func (c *UploadController) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	// Get file from form
	file, header, err := r.FormFile("image")
	if err != nil {
		respondError(w, http.StatusBadRequest, "Image file is required")
		return
	}
	defer file.Close()

	// Upload image
	filename, err := c.uploadService.UploadImage(file, header)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get image URL
	imageURL := c.uploadService.GetImageURL(filename)

	data := map[string]interface{}{
		"filename": filename,
		"url":      imageURL,
	}

	respondSuccess(w, http.StatusOK, "Image uploaded successfully", data)
}

// ServeImage serves uploaded images
func (c *UploadController) ServeImage(w http.ResponseWriter, r *http.Request) {
	// Get filename from URL path
	filename := strings.TrimPrefix(r.URL.Path, "/uploads/books/")
	if filename == "" {
		http.Error(w, "Filename required", http.StatusBadRequest)
		return
	}

	// Validate filename (prevent directory traversal)
	if err := service.ValidateImagePath(filename); err != nil {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(c.uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	// Set cache headers for better performance
	w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year

	// Serve the file
	http.ServeFile(w, r, filePath)
}

// DeleteImage handles image deletion (for cleanup when book is deleted)
func (c *UploadController) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.Query().Get("filename")
	if filename == "" {
		respondError(w, http.StatusBadRequest, "Filename is required")
		return
	}

	if err := service.ValidateImagePath(filename); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid filename")
		return
	}

	err := c.uploadService.DeleteImage(filename)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Image deleted successfully", nil)
}
