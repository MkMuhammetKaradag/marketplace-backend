package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {

	if _, err := db.Exec(createUserRolesEnum); err != nil {
		return fmt.Errorf("failed to create user roles enum: %w", err)
	}

	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Database migrated")
	return nil
}
