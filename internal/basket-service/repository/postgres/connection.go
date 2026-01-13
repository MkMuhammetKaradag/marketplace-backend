// internal/basket-service/repository/postgres/connection.go
package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"marketplace/internal/basket-service/config"
)

func newPostgresDB(cfg config.Config) (*sql.DB, error) {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
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
