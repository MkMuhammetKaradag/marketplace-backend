package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {
	if _, err := db.Exec(createSellerStatusEnum); err != nil {
		return fmt.Errorf("failed to create seller status enum: %w", err)
	}

	if _, err := db.Exec(createSellersTable); err != nil {
		return fmt.Errorf("failed to create sellers table: %w", err)
	}
	if _, err := db.Exec(createSellerStatusHistoryTable); err != nil {
		return fmt.Errorf("failed to create seller status history table: %w", err)
	}

	log.Println("Database migrated")
	return nil
}
