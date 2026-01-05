package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const UPSERT_USER_PREFERENCE = `
    INSERT INTO user_preferences (user_id, interest_vector, last_interaction_at, updated_at)
    VALUES ($1, $2, NOW(), NOW())
    ON CONFLICT (user_id) DO UPDATE SET
        -- Çarpma yerine doğrudan vektör toplama kullanıyoruz
        -- Bu işlem genellikle tüm pgvector sürümlerinde sorunsuz çalışır
        interest_vector = user_preferences.interest_vector + EXCLUDED.interest_vector,
        last_interaction_at = NOW(),
        updated_at = NOW();
`

func (r *Repository) TrackProductView(ctx context.Context, userID uuid.UUID, productEmbedding []float32) error {

	bytes, err := json.Marshal(productEmbedding)
	if err != nil {
		return fmt.Errorf("failed to marshal product embedding: %w", err)
	}

	_, err = r.db.ExecContext(ctx, UPSERT_USER_PREFERENCE, userID, string(bytes))
	if err != nil {
		return fmt.Errorf("failed to upsert user preference: %w", err)
	}

	return nil
}
