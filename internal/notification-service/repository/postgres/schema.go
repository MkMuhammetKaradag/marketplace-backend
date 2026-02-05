package postgres

const (
	createLocalUsersTable = `
		CREATE TABLE IF NOT EXISTS local_users (
            id UUID PRIMARY KEY,
            username VARCHAR(100), 
            email VARCHAR(100),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )`
)
