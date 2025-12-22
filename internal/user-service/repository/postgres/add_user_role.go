package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const query = `
        INSERT INTO user_roles (user_id, role_id)
        SELECT $1, id FROM roles WHERE name = $2
        ON CONFLICT (user_id, role_id) DO NOTHING`

func (r *Repository) AddUserRole(ctx context.Context, userID uuid.UUID, roleName string) error {
	result, err := r.db.ExecContext(ctx, query, userID, roleName)
	if err != nil {
		return fmt.Errorf("failed to add user role: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		fmt.Println("User role not found")
	}

	return nil
}
