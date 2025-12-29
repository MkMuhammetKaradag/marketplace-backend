package postgres

const (
	createSellerStatusEnum = `
        DO $$
        BEGIN
            IF NOT EXISTS (
                SELECT 1 FROM pg_type WHERE typname = 'seller_status'
            ) THEN
                CREATE TYPE seller_status AS ENUM ('pending', 'approved', 'rejected');
            END IF;
        END$$;
    `
	createLocalSellersTable = `
		CREATE TABLE IF NOT EXISTS local_sellers (
			seller_id UUID PRIMARY KEY,
			status seller_status, -- Approved mu?
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()

		)`
)
