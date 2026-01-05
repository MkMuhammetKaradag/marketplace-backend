package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const ADD_USER_INTERACTION = `
    INSERT INTO user_product_interactions (user_id, product_id, interaction_type, weight, created_at)
    VALUES ($1, $2, $3, $4, NOW())
`

func (r *Repository) AddInteraction(ctx context.Context, userID uuid.UUID, productID uuid.UUID, interactionType string) error {
	// Etkileşim tipine göre ağırlık belirleyelim
	var weight float64
	switch interactionType {
	case "view":
		weight = 1.0
	case "like":
		weight = 3.0
	case "purchase":
		weight = 5.0
	default:
		weight = 1.0
	}

	_, err := r.db.ExecContext(ctx, ADD_USER_INTERACTION, userID, productID, interactionType, weight)
	if err != nil {
		return fmt.Errorf("failed to add interaction: %w", err)
	}

	return nil
}
