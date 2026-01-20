package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {

	if _, err := db.Exec(createOrdersTable); err != nil {
		return fmt.Errorf("failed to create orders table: %w", err)
	}

	if _, err := db.Exec(createOrderItemsTable); err != nil {
		return fmt.Errorf("failed to create order_items table: %w", err)
	}

	// if _, err := db.Exec(createIndex); err != nil {
	// 	return fmt.Errorf("failed to create index: %w", err)
	// }

	log.Println("Database migration completed successfully")
	return nil
}
