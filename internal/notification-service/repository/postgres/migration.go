package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {

	if _, err := db.Exec(createLocalUsersTable); err != nil {
		return fmt.Errorf("failed to create local_users table: %w", err)
	}

	log.Println("Database migration completed successfully")
	return nil
}
