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

func (r *Repository) RejectSeller(ctx context.Context, sellerId string, rejectedBy string, rejectionReason string) (string, error) {
	row := r.db.QueryRowContext(ctx, rejectSellerQuery, rejectedBy, rejectionReason, sellerId)
	var userId string
	if err := row.Scan(&userId); err != nil {
		return "", err
	}
	return userId, nil
}
