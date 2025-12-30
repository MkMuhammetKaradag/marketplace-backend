package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {

	if _, err := db.Exec(createSellerStatusEnum); err != nil {
		return fmt.Errorf("failed to create seller_status enum: %w", err)
	}
	if _, err := db.Exec(createLocalSellersTable); err != nil {
		return fmt.Errorf("failed to create local_sellers table: %w", err)
	}
	if _, err := db.Exec(createProductStatusEnum); err != nil {
		return fmt.Errorf("failed to create product_status enum: %w", err)
	}
	if _, err := db.Exec(createProductsTable); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
