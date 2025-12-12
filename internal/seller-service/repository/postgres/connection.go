package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"marketplace/internal/seller-service/config"
)

func newPostgresDB(cfg config.Config) (*sql.DB, error) {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DB,
	)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Connected to PostgreSQL")
	return db, nil
}
