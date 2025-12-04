package postgres

const (
	createUsersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    		username VARCHAR(50) NOT NULL UNIQUE,
    		email VARCHAR(100) NOT NULL UNIQUE,
    		password TEXT NOT NULL,
    		is_active BOOLEAN DEFAULT false,
    		is_email_verified BOOLEAN DEFAULT false,
			activation_code VARCHAR(6),
			activation_id UUID  DEFAULT gen_random_uuid(),
			activation_expiry TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			failed_login_attempts INT DEFAULT 0,
			account_locked BOOLEAN DEFAULT false,
			lock_until TIMESTAMP WITH TIME ZONE,
			last_login TIMESTAMP WITH TIME ZONE,
			CONSTRAINT check_activation_code CHECK (
				activation_code ~ '^[0-9]{6}$'
			)
		)`
)
