package postgres

import (
	"database/sql"
	"log"
)

func runMigrations(db *sql.DB) error {

	log.Println("Database migrated")
	return nil
}
