package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {

	if _, err := db.Exec(createExtension); err != nil {
		return fmt.Errorf("failed to create extension: %w", err)
	}
	if _, err := db.Exec(createSellerStatusEnum); err != nil {
		return fmt.Errorf("failed to create seller_status enum: %w", err)
	}
	if _, err := db.Exec(createLocalUsersTable); err != nil {
		return fmt.Errorf("failed to create local_users table: %w", err)
	}
	if _, err := db.Exec(createLocalSellersTable); err != nil {
		return fmt.Errorf("failed to create local_sellers table: %w", err)
	}
	if _, err := db.Exec(createCategoriesTable); err != nil {
		return fmt.Errorf("failed to create categories table: %w", err)
	}
	if _, err := db.Exec(createProductStatusEnum); err != nil {
		return fmt.Errorf("failed to create product_status enum: %w", err)
	}
	if _, err := db.Exec(createProductsTable); err != nil {
		return fmt.Errorf("failed to create products table: %w", err)
	}
	if _, err := db.Exec(createProductImagesTable); err != nil {
		return fmt.Errorf("failed to create product_images table: %w", err)
	}
	if _, err := db.Exec(createUserPreferencesTable); err != nil {
		return fmt.Errorf("failed to create user_preferences table: %w", err)
	}
	if _, err := db.Exec(createUserProductInteractionsTable); err != nil {
		return fmt.Errorf("failed to create user_product_interactions table: %w", err)
	}
	if _, err := db.Exec(createFavoriteTable); err != nil {
		return fmt.Errorf("failed to create favorites table: %w", err)
	}
	if _, err := db.Exec(createIndex); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	if _, err := db.Exec(createCleanupProductFunction); err != nil {
		return fmt.Errorf("failed to create cleanup_product function: %w", err)
	}
	log.Println("Database migration completed successfully")
	return nil
}
