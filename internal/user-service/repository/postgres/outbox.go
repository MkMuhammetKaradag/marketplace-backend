package postgres

import (
	"context"
	"marketplace/internal/user-service/domain"

	"github.com/google/uuid"
)

func (r *Repository) GetPendingOutboxMessages(ctx context.Context, limit int) ([]domain.OutboxMessage, error) {
	query := `
        SELECT id, payload, topic 
        FROM outbox_messages 
        WHERE status = 'PENDING' 
        ORDER BY created_at ASC 
        LIMIT $1`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.OutboxMessage
	for rows.Next() {
		var msg domain.OutboxMessage
		if err := rows.Scan(&msg.ID, &msg.Payload, &msg.Topic); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *Repository) MarkOutboxAsProcessed(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE outbox_messages SET status = 'PROCESSED', updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
