package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const REMOVE_USER_INTERACTION = `
    DELETE FROM user_product_interactions 
    WHERE user_id = $1 AND product_id = $2 AND interaction_type = $3
`

func (r *Repository) RemoveInteraction(ctx context.Context, userID uuid.UUID, productID uuid.UUID, interactionType string) error {
	_, err := r.db.ExecContext(ctx, REMOVE_USER_INTERACTION, userID, productID, interactionType)
	if err != nil {
		return fmt.Errorf("failed to remove interaction: %w", err)
	}
	return nil
}
