package router

import (
	"net/http"

	"github.com/LanangDepok/ebook-store/controller"
	"github.com/LanangDepok/ebook-store/middleware"
)

type Router struct {
	authController   *controller.AuthController
	bookController   *controller.BookController
	cartController   *controller.CartController
	orderController  *controller.OrderController
	uploadController *controller.UploadController
	authMiddleware   *middleware.AuthMiddleware
}

func NewRouter(
	authController *controller.AuthController,
	bookController *controller.BookController,
	cartController *controller.CartController,
	orderController *controller.OrderController,
	uploadController *controller.UploadController,
	authMiddleware *middleware.AuthMiddleware,
) *Router {
	return &Router{
		authController:   authController,
		bookController:   bookController,
		cartController:   cartController,
		orderController:  orderController,
		uploadController: uploadController,
		authMiddleware:   authMiddleware,
	}
}

func (router *Router) Setup() *http.ServeMux {
	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("/api/auth/register", methodHandler("POST", router.authController.Register))
	mux.HandleFunc("/api/auth/login", methodHandler("POST", router.authController.Login))
	mux.HandleFunc("/api/auth/logout", methodHandler("POST", router.authMiddleware.RequireAuth(router.authController.Logout)))

	// Book routes
	mux.HandleFunc("/api/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			router.bookController.GetAllBooks(w, r)
		case "POST":
			router.authMiddleware.RequireAdmin(router.bookController.CreateBook)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/books/detail", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			router.bookController.GetBookByID(w, r)
		case "PUT":
			router.authMiddleware.RequireAdmin(router.bookController.UpdateBook)(w, r)
		case "DELETE":
			router.authMiddleware.RequireAdmin(router.bookController.DeleteBook)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Cart routes
	mux.HandleFunc("/api/cart", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			router.authMiddleware.RequireAuth(router.cartController.GetCart)(w, r)
		case "POST":
			router.authMiddleware.RequireAuth(router.cartController.AddToCart)(w, r)
		case "DELETE":
			router.authMiddleware.RequireAuth(router.cartController.ClearCart)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/cart/item", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "PUT":
			router.authMiddleware.RequireAuth(router.cartController.UpdateCartItem)(w, r)
		case "DELETE":
			router.authMiddleware.RequireAuth(router.cartController.RemoveFromCart)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Order routes
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			router.authMiddleware.RequireAuth(router.orderController.GetUserOrders)(w, r)
		case "POST":
			router.authMiddleware.RequireAuth(router.orderController.CreateOrder)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/orders/detail", methodHandler("GET", router.authMiddleware.RequireAuth(router.orderController.GetOrderDetail)))

	// Health check
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy"}`))
	})

	return mux
}

func methodHandler(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, r)
	}
}
