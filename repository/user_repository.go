package repository

import (
	"database/sql"
	"fmt"

	"github.com/LanangDepok/ebook-store/entity"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByID(id int) (*entity.User, error)
	UsernameExists(username string) (bool, error)
	EmailExists(email string) (bool, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	query := `
		INSERT INTO users (username, password, email, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(query, user.Username, user.Password, user.Email, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) FindByUsername(username string) (*entity.User, error) {
	query := `
		SELECT id, username, password, email, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	user := &entity.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	query := `
		SELECT id, username, password, email, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	user := &entity.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByID(id int) (*entity.User, error) {
	query := `
		SELECT id, username, password, email, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &entity.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email,
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) UsernameExists(username string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&count)
	return count > 0, err
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(&count)
	return count > 0, err
}

type SessionRepository interface {
	Create(session *entity.Session) error
	FindByID(id string) (*entity.Session, error)
	Delete(id string) error
	DeleteExpired() error
}

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *entity.Session) error {
	query := `
		INSERT INTO sessions (id, user_id, username, created_at, expires_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, session.ID, session.UserID, session.Username,
		session.CreatedAt, session.ExpiresAt)
	return err
}

func (r *sessionRepository) FindByID(id string) (*entity.Session, error) {
	query := `
		SELECT id, user_id, username, created_at, expires_at
		FROM sessions
		WHERE id = $1
	`
	session := &entity.Session{}
	err := r.db.QueryRow(query, id).Scan(
		&session.ID, &session.UserID, &session.Username,
		&session.CreatedAt, &session.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (r *sessionRepository) Delete(id string) error {
	query := `DELETE FROM sessions WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("session not found")
	}
	return nil
}

func (r *sessionRepository) DeleteExpired() error {
	query := `DELETE FROM sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(query)
	return err
}
