package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LanangDepok/ebook-store/middleware"
	"github.com/LanangDepok/ebook-store/model"
	"github.com/LanangDepok/ebook-store/service"
)

type CartController struct {
	cartService service.CartService
}

func NewCartController(cartService service.CartService) *CartController {
	return &CartController{cartService: cartService}
}

func (c *CartController) AddToCart(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req model.AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.BookID <= 0 || req.Jumlah <= 0 {
		respondError(w, http.StatusBadRequest, "Invalid cart data")
		return
	}

	err := c.cartService.AddToCart(user.ID, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Item added to cart successfully", nil)
}

func (c *CartController) GetCart(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	items, total, err := c.cartService.GetCart(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := map[string]interface{}{
		"items": items,
		"total": total,
	}

	respondSuccess(w, http.StatusOK, "Cart retrieved successfully", data)
}

func (c *CartController) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Cart item ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	var req model.UpdateCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Jumlah <= 0 {
		respondError(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	err = c.cartService.UpdateCartItem(id, req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Cart item updated successfully", nil)
}

func (c *CartController) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "Cart item ID is required")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid cart item ID")
		return
	}

	err = c.cartService.RemoveFromCart(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Item removed from cart successfully", nil)
}

func (c *CartController) ClearCart(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r.Context())
	if user == nil {
		respondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := c.cartService.ClearCart(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Cart cleared successfully", nil)
}
