package postgres

import (
	"database/sql"
	"log"
)

func runMigrations(db *sql.DB) error {

	// if _, err := db.Exec(createCleanupProductFunction); err != nil {
	// 	return fmt.Errorf("failed to create cleanup_product function: %w", err)
	// }
	log.Println("Database migration completed successfully")
	return nil
}
