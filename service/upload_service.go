package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UploadService interface {
	UploadImage(file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteImage(filename string) error
	GetImageURL(filename string) string
}

type uploadService struct {
	uploadDir string
	baseURL   string
}

func NewUploadService(uploadDir, baseURL string) UploadService {
	return &uploadService{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (s *uploadService) UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		return "", fmt.Errorf("file size exceeds 5MB limit")
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !s.isValidImageType(contentType) {
		return "", fmt.Errorf("invalid file type. Only JPEG, PNG, GIF, and WebP are allowed")
	}

	// Create upload directory if not exists
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %v", err)
	}

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), generateRandomString(8), ext)
	filePath := filepath.Join(s.uploadDir, filename)

	// Create destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy uploaded file to destination
	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath) // Clean up on error
		return "", fmt.Errorf("failed to save file: %v", err)
	}

	return filename, nil
}

func (s *uploadService) DeleteImage(filename string) error {
	if filename == "" {
		return nil
	}

	filePath := filepath.Join(s.uploadDir, filename)
	if err := os.Remove(filePath); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete image: %v", err)
		}
	}

	return nil
}

func (s *uploadService) GetImageURL(filename string) string {
	if filename == "" {
		return ""
	}
	return fmt.Sprintf("%s/uploads/books/%s", s.baseURL, filename)
}

func (s *uploadService) isValidImageType(contentType string) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	return validTypes[contentType]
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(b)
}

// ValidateImagePath prevents directory traversal attacks
func ValidateImagePath(filename string) error {
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("invalid filename")
	}
	return nil
}
