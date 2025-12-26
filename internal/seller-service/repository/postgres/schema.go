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
	createSellersTable = `
        CREATE TABLE IF NOT EXISTS sellers (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            
            user_id UUID NOT NULL UNIQUE, 
            
            store_name VARCHAR(100) NOT NULL UNIQUE,
            store_slug VARCHAR(100) NOT NULL UNIQUE, 
            store_logo_url VARCHAR(255),
            store_banner_url VARCHAR(255),
            store_description TEXT,
            rating NUMERIC(2, 1) DEFAULT 0.0, 
            total_sales INT DEFAULT 0,
            legal_business_name VARCHAR(255) NOT NULL,
            tax_number VARCHAR(20) UNIQUE, 
            tax_office VARCHAR(100) NOT NULL,
            is_approved BOOLEAN DEFAULT false, 

            phone_number VARCHAR(20) NOT NULL,
            email VARCHAR(100) UNIQUE, 
            
            address_line_1 VARCHAR(255) NOT NULL, 
            city VARCHAR(100) NOT NULL,
            country VARCHAR(50) NOT NULL,
            default_shipping_country VARCHAR(50),	

            bank_account_iban VARCHAR(34) UNIQUE, 
            bank_account_holder_name VARCHAR(255) NOT NULL,
            bank_account_bic VARCHAR(11), 
            status seller_status DEFAULT 'pending',
            approved_at TIMESTAMP WITH TIME ZONE,
            rejected_at TIMESTAMP WITH TIME ZONE,
            rejection_reason TEXT,
            approved_by UUID,
            rejected_by UUID,

    
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            
           	
            CONSTRAINT check_tax_number_format CHECK (
                 tax_number ~ '^[0-9]{10,20}$' 	
            )
        )`

	createSellerStatusHistoryTable = `
        CREATE TABLE IF NOT EXISTS seller_status_history (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            seller_id UUID NOT NULL REFERENCES sellers(id) ON DELETE CASCADE,
            status seller_status NOT NULL,
            reason TEXT,
            changed_by UUID,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )`
)
