package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Database struct {
	DB *sql.DB
}

func NewDatabase() *Database {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, name, sslmode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	return &Database{DB: db}
}

func (d *Database) Close() error {
	return d.DB.Close()
}

func (d *Database) Migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			role VARCHAR(20) DEFAULT 'user',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id VARCHAR(64) PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			username VARCHAR(50) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			expires_at TIMESTAMP NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			nama_barang TEXT NOT NULL,
			stok INTEGER DEFAULT 0,
			terjual INTEGER DEFAULT 0,
			harga INTEGER DEFAULT 0,
			keterangan TEXT,
			gambar_buku TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS carts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
			jumlah INTEGER NOT NULL DEFAULT 1,
			harga INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(user_id, book_id)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			total_harga INTEGER NOT NULL DEFAULT 0,
			status VARCHAR(50) DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id SERIAL PRIMARY KEY,
			order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
			book_id INTEGER NOT NULL REFERENCES books(id) ON DELETE CASCADE,
			jumlah INTEGER NOT NULL DEFAULT 1,
			harga INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at)`,
		`CREATE INDEX IF NOT EXISTS idx_carts_user_id ON carts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
	}

	for _, migration := range migrations {
		if _, err := d.DB.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %v", err)
		}
	}

	// Insert default users with hashed passwords
	adminPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)

	insertUsers := `INSERT INTO users (username, password, email, role) VALUES 
		 ($1, $2, 'admin@book.com', 'admin'),
		 ($3, $4, 'user@book.com', 'user')
		 ON CONFLICT (username) DO NOTHING`

	if _, err := d.DB.Exec(insertUsers, "admin", string(adminPassword), "user", string(userPassword)); err != nil {
		return fmt.Errorf("failed to insert default users: %v", err)
	}

	log.Println("Database migration completed")
	return nil
}
