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

	createProductStatusEnum = `
        DO $$ 
            BEGIN
                CREATE TYPE product_status AS ENUM ('active', 'inactive', 'out_of_stock', 'deleted');
            EXCEPTION
                WHEN duplicate_object THEN null;
        END $$;
    `
	createLocalSellersTable = `
		CREATE TABLE IF NOT EXISTS local_sellers (
			seller_id UUID PRIMARY KEY,
			status seller_status, -- Approved mu?
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()

		)`

	createProductsTable = `
		CREATE TABLE IF NOT EXISTS products (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            seller_id UUID NOT NULL,      
            name VARCHAR(255) NOT NULL,
            description TEXT,
            price DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
            stock_count INTEGER NOT NULL DEFAULT 0,
            status product_status DEFAULT 'inactive',
            attributes JSONB DEFAULT '{}',
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

            
            
            CONSTRAINT fk_seller
                FOREIGN KEY(seller_id) 
                REFERENCES local_sellers(seller_id)
                ON DELETE CASCADE 
        );`
)
