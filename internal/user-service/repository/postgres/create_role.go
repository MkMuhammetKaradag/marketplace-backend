package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const createRoleQuery = `
    INSERT INTO roles (
        created_by, name, permissions, color, position, 
        is_mentionable, is_hoisted, is_managed
    ) VALUES (
        NULLIF($1, '00000000-0000-0000-0000-000000000000'::uuid), -- Sıfırsa NULL basar
        $2, $3, $4, $5, $6, $7, $8
    )
    RETURNING id`

func (r *Repository) CreateRole(ctx context.Context, createdBy uuid.UUID, name string, permissions int64) (uuid.UUID, error) {

	var roleID uuid.UUID

	color := "#B9BBBE"
	position := 0
	isMentionable := true
	isHoisted := false
	isManaged := false

	err := r.db.QueryRowContext(ctx, createRoleQuery,
		createdBy,
		name,
		permissions,
		color,
		position,
		isMentionable,
		isHoisted,
		isManaged,
	).Scan(&roleID)

	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create role: %w", err)
	}

	return roleID, nil
}
