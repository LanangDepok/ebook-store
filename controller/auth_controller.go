package controller

import (
	"encoding/json"
	"net/http"

	"github.com/LanangDepok/ebook-store/model"
	"github.com/LanangDepok/ebook-store/service"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Username == "" || req.Password == "" || req.Email == "" {
		respondError(w, http.StatusBadRequest, "All fields are required")
		return
	}

	if len(req.Password) < 6 {
		respondError(w, http.StatusBadRequest, "Password must be at least 6 characters")
		return
	}

	err := c.authService.Register(req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, "User registered successfully", nil)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		respondError(w, http.StatusBadRequest, "Username and password are required")
		return
	}

	response, err := c.authService.Login(req)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Login successful", response)
}

func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	token := extractToken(r)
	if token == "" {
		respondError(w, http.StatusBadRequest, "Missing token")
		return
	}

	err := c.authService.Logout(token)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to logout")
		return
	}

	respondSuccess(w, http.StatusOK, "Logout successful", nil)
}

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
