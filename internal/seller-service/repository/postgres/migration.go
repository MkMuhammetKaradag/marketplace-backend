package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec(createSellersTable); err != nil {
		return fmt.Errorf("failed to create sellers table: %w", err)
	}

	log.Println("Database migrated")
	return nil
}
