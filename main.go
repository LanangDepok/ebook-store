package main

import (
	"log"
	"net/http"
	"os"

	"github.com/LanangDepok/ebook-store/config"
	"github.com/LanangDepok/ebook-store/controller"
	"github.com/LanangDepok/ebook-store/middleware"
	"github.com/LanangDepok/ebook-store/repository"
	"github.com/LanangDepok/ebook-store/router"
	"github.com/LanangDepok/ebook-store/service"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Create uploads directory
	uploadDir := "uploads/books"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	log.Printf("Upload directory ready: %s", uploadDir)

	// Initialize database
	db := config.NewDatabase()
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Get base URL for images
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		baseURL = "http://localhost:" + port
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)
	sessionRepo := repository.NewSessionRepository(db.DB)
	bookRepo := repository.NewBookRepository(db.DB)
	cartRepo := repository.NewCartRepository(db.DB)
	orderRepo := repository.NewOrderRepository(db.DB)

	// Initialize services
	authService := service.NewAuthService(userRepo, sessionRepo)
	uploadService := service.NewUploadService(uploadDir, baseURL)
	bookService := service.NewBookService(bookRepo)
	cartService := service.NewCartService(cartRepo, bookRepo)
	orderService := service.NewOrderService(orderRepo, cartRepo, bookRepo, db.DB)

	// Initialize controllers
	authController := controller.NewAuthController(authService)
	bookController := controller.NewBookController(bookService, uploadService)
	cartController := controller.NewCartController(cartService, uploadService)
	orderController := controller.NewOrderController(orderService)
	uploadController := controller.NewUploadController(uploadService, uploadDir)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(db.DB)

	// Setup router
	appRouter := router.NewRouter(
		authController,
		bookController,
		cartController,
		orderController,
		uploadController,
		authMiddleware,
	)

	mux := appRouter.Setup()

	// Apply middleware
	handler := middleware.CORS(mux)
	handler = middleware.ContentTypeJSON(handler)

	// Get port from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Println("========================================")
	log.Printf("🚀 Server starting on http://localhost:%s", port)
	log.Println("========================================")
	log.Println("Available endpoints:")
	log.Println("  Auth:")
	log.Println("    POST   /api/auth/register")
	log.Println("    POST   /api/auth/login")
	log.Println("    POST   /api/auth/logout")
	log.Println("  Books:")
	log.Println("    GET    /api/books")
	log.Println("    POST   /api/books (admin only)")
	log.Println("    GET    /api/books/detail?id=1")
	log.Println("    PUT    /api/books/detail?id=1 (admin only)")
	log.Println("    DELETE /api/books/detail?id=1 (admin only)")
	log.Println("  Cart:")
	log.Println("    GET    /api/cart")
	log.Println("    POST   /api/cart")
	log.Println("    PUT    /api/cart/item?id=1")
	log.Println("    DELETE /api/cart/item?id=1")
	log.Println("    DELETE /api/cart (clear cart)")
	log.Println("  Orders:")
	log.Println("    GET    /api/orders")
	log.Println("    POST   /api/orders")
	log.Println("    GET    /api/orders/detail?id=1")
	log.Println("  Upload:")
	log.Println("    POST   /api/upload/image (admin only)")
	log.Println("  Static:")
	log.Println("    GET    /uploads/books/{filename}")
	log.Println("  Health:")
	log.Println("    GET    /api/health")
	log.Println("========================================")
	log.Println("Default users:")
	log.Println("  Admin: username=admin, password=admin123")
	log.Println("  User:  username=user, password=user123")
	log.Println("========================================")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
