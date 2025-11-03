package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/LanangDepok/ebook-store/entity"
	"github.com/LanangDepok/ebook-store/model"
	"github.com/LanangDepok/ebook-store/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req model.RegisterRequest) error
	Login(req model.LoginRequest) (*model.LoginResponse, error)
	Logout(token string) error
	ValidateToken(token string) (*entity.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (s *authService) Register(req model.RegisterRequest) error {
	// Check if username exists
	exists, err := s.userRepo.UsernameExists(req.Username)
	if err != nil {
		return fmt.Errorf("failed to check username: %v", err)
	}
	if exists {
		return fmt.Errorf("username already exists")
	}

	// Check if email exists
	exists, err = s.userRepo.EmailExists(req.Email)
	if err != nil {
		return fmt.Errorf("failed to check email: %v", err)
	}
	if exists {
		return fmt.Errorf("email already exists")
	}

	// Hash password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Create user
	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Role:     "user",
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Verify password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate session token
	token := generateToken()
	session := &entity.Session{
		ID:        token,
		UserID:    user.ID,
		Username:  user.Username,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = s.sessionRepo.Create(session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &model.LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	}, nil
}

func (s *authService) Logout(token string) error {
	return s.sessionRepo.Delete(token)
}

func (s *authService) ValidateToken(token string) (*entity.User, error) {
	session, err := s.sessionRepo.FindByID(token)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	if time.Now().After(session.ExpiresAt) {
		s.sessionRepo.Delete(token)
		return nil, fmt.Errorf("session expired")
	}

	return s.userRepo.FindByID(session.UserID)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
