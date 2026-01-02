package postgres

import (
	"context"
	"marketplace/internal/product-service/domain"

	"github.com/google/uuid"
)

const CREATE_CATEGORY = `
    INSERT INTO categories (parent_id, name, slug, description)
    VALUES ($1, $2, $3, $4)
    RETURNING id`

func (r *Repository) CreateCategory(ctx context.Context, c *domain.Category) error {
	var parentID interface{}
	if c.ParentID == uuid.Nil {
		parentID = nil
	} else {
		parentID = c.ParentID
	}

	return r.db.QueryRowContext(ctx, CREATE_CATEGORY,
		parentID, c.Name, c.Slug, c.Description,
	).Scan(&c.ID)

}
