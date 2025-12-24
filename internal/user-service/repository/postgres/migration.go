package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

func runMigrations(db *sql.DB) error {
	// 1. Sıra: Users Tablosu (Roles tablosu buna referans veriyor)
	if _, err := db.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// 2. Sıra: Roles Tablosu
	if _, err := db.Exec(createRolesTable); err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	// 3. Sıra: User_Roles (İlişki Tablosu)
	if _, err := db.Exec(createdUserRolesTable); err != nil {
		return fmt.Errorf("failed to create user_roles table: %w", err)
	}

	// 4. Sıra: Varsayılan Rolleri Ekleme
	// Not: ON CONFLICT DO NOTHING ekleyerek her restartta hata almayı önlüyoruz
	if _, err := db.Exec(createDefaultRoles); err != nil {
		return fmt.Errorf("failed to insert default roles: %w", err)
	}

	// 5. Sıra: Forgot Passwords Tablosu
	if _, err := db.Exec(createForgotPasswordsTable); err != nil {
		return fmt.Errorf("failed to create forgot_passwords table: %w", err)
	}
	
	log.Println("Database migration completed successfully")
	return nil
}
