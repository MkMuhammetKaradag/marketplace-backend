package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

const UPDATE_PRODUCT_EMBEDDING = `
    UPDATE products 
    SET embedding = $1 
    WHERE id = $2`

func (r *Repository) UpdateProductEmbedding(ctx context.Context, id uuid.UUID, embedding []float32) error {
	// pgvector, veriyi "[0.1, 0.2, 0.3]" formatında bir string olarak bekler.
	// json.Marshal, float slice'ını tam olarak bu formata (virgüllü ve köşeli parantezli) çevirir.
	embeddingBytes, err := json.Marshal(embedding)
	if err != nil {
		return fmt.Errorf("failed to marshal embedding: %w", err)
	}

	// string(embeddingBytes) artık "[0.123,0.456,...]" şeklindedir.
	_, err = r.db.ExecContext(ctx, UPDATE_PRODUCT_EMBEDDING, string(embeddingBytes), id)
	if err != nil {
		return fmt.Errorf("failed to update product embedding: %w", err)
	}
	fmt.Println("Product embedding updated successfully")

	return nil
}
