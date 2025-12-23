package postgres

import "context"

const rejectSellerQuery = `
    UPDATE sellers
    SET status = 'rejected',
        rejected_at = NOW(),
        rejected_by = $1,
		rejection_reason = $2
    WHERE id = $3 AND status = 'pending'
    RETURNING user_id;
`

const rejectSellerHistoryQuery = `
    INSERT INTO seller_status_history (seller_id, status, reason, changed_by) 
        VALUES ($1, 'rejected', $2, $3);
`

func (r *Repository) RejectSeller(ctx context.Context, sellerId string, rejectedBy string, rejectionReason string) (string, error) {
	row := r.db.QueryRowContext(ctx, rejectSellerQuery, rejectedBy, rejectionReason, sellerId)
	var userId string
	if err := row.Scan(&userId); err != nil {
		return "", err
	}
	_, err := r.db.ExecContext(ctx, rejectSellerHistoryQuery, sellerId, rejectionReason, rejectedBy)
	if err != nil {
		return "", err
	}
	return userId, nil
}
