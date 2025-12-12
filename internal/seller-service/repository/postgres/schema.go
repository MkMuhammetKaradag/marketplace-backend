package postgres

const (
	createSellersTable = `
		CREATE TABLE IF NOT EXISTS sellers (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			
			user_id UUID NOT NULL UNIQUE, 
			store_name VARCHAR(100) NOT NULL UNIQUE,
			store_slug VARCHAR(100) NOT NULL UNIQUE, 
			tax_number VARCHAR(20) UNIQUE, 
			is_approved BOOLEAN DEFAULT false, 
			default_shipping_country VARCHAR(50),
			store_description TEXT,
			
			bank_account_iban VARCHAR(34), 
			
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			
		
			rating NUMERIC(2, 1) DEFAULT 0.0, 
			total_sales INT DEFAULT 0
		)`
)
