package middleware

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/LanangDepok/ebook-store/entity"
	"github.com/LanangDepok/ebook-store/model"
)

type contextKey string

const UserContextKey contextKey = "user"

type AuthMiddleware struct {
	db *sql.DB
}

func NewAuthMiddleware(db *sql.DB) *AuthMiddleware {
	return &AuthMiddleware{db: db}
}

func (m *AuthMiddleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := m.extractToken(r)
		if token == "" {
			m.unauthorizedResponse(w, "Missing authentication token")
			return
		}

		user, err := m.validateToken(token)
		if err != nil {
			m.unauthorizedResponse(w, "Invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		user := GetUserFromContext(r.Context())
		if user == nil || user.Role != "admin" {
			m.forbiddenResponse(w, "Admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (m *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func (m *AuthMiddleware) validateToken(token string) (*entity.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.role
		FROM users u
		JOIN sessions s ON u.id = s.user_id
		WHERE s.id = $1 AND s.expires_at > $2
	`

	var user entity.User
	err := m.db.QueryRow(query, token, time.Now()).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *AuthMiddleware) unauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(model.Response{
		Status:  "error",
		Message: message,
	})
}

func (m *AuthMiddleware) forbiddenResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(model.Response{
		Status:  "error",
		Message: message,
	})
}

func GetUserFromContext(ctx context.Context) *entity.User {
	user, ok := ctx.Value(UserContextKey).(*entity.User)
	if !ok {
		return nil
	}
	return user
}
