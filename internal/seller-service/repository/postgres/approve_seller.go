package postgres

import "context"

const approveSellerQuery = `
    UPDATE sellers
    SET status = 'approved',
        approved_at = NOW(),
        approved_by = $1
    WHERE id = $2 AND status = 'pending'
    RETURNING user_id;
`

func (r *Repository) ApproveSeller(ctx context.Context, sellerId string, approvedBy string) (string, error) {
	row := r.db.QueryRowContext(ctx, approveSellerQuery, approvedBy, sellerId)
	var userId string
	if err := row.Scan(&userId); err != nil {
		return "", err
	}
	return userId, nil
}
