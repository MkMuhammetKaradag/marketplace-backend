package postgres

const (
	createExtension = `
        CREATE EXTENSION IF NOT EXISTS vector
    `

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
            user_id UUID NOT NULL UNIQUE,
			status seller_status,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()

		)`
	createLocalUsersTable = `
		CREATE TABLE IF NOT EXISTS local_users (
            id UUID PRIMARY KEY,
            username VARCHAR(100), 
            email VARCHAR(100),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )`

	createProductsTable = `
		CREATE TABLE IF NOT EXISTS products (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            seller_id UUID NOT NULL,      
            name VARCHAR(255) NOT NULL,
            category_id UUID REFERENCES categories(id) NOT NULL,
            description TEXT,
            price DECIMAL(12, 2) NOT NULL DEFAULT 0.00,
            stock_count INTEGER NOT NULL DEFAULT 0,
            status product_status DEFAULT 'inactive',
            attributes JSONB DEFAULT '{}',
            embedding vector(768),
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),

            
            
            CONSTRAINT fk_seller
                FOREIGN KEY(seller_id) 
                REFERENCES local_sellers(seller_id)
                ON DELETE CASCADE 
        );`

	createProductImagesTable = `
		CREATE TABLE IF NOT EXISTS product_images (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			product_id UUID NOT NULL,
			image_url TEXT NOT NULL,
            is_main BOOLEAN DEFAULT FALSE,
            sort_order INTEGER DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            deleted_at TIMESTAMP WITH TIME ZONE, 

			CONSTRAINT fk_product
			    FOREIGN KEY(product_id) 
			    REFERENCES products(id)
			    ON DELETE CASCADE 
		);`

	createCategoriesTable = `
        CREATE TABLE IF NOT EXISTS categories (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            parent_id UUID REFERENCES categories(id) ON DELETE CASCADE, -- Alt kategori desteği
            name VARCHAR(100) NOT NULL,
            slug VARCHAR(100) NOT NULL UNIQUE, -- URL dostu isim (örn: bilgisayar-laptop)
            description TEXT,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )`

	createUserPreferencesTable = `
        CREATE TABLE IF NOT EXISTS user_preferences (
            user_id UUID PRIMARY KEY, -- Harici user tablosuna referans
            interest_vector vector(768),
            last_interaction_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )
        `
	createUserProductInteractionsTable = `
        CREATE TABLE IF NOT EXISTS user_product_interactions (
            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
            user_id UUID NOT NULL,
            product_id UUID REFERENCES products(id) ON DELETE CASCADE,
            interaction_type VARCHAR(20), -- 'view', 'like', 'purchase'
            weight FLOAT DEFAULT 1.0,     -- Satın alma: 5.0, İzleme: 1.0 gibi
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
        )  
    `

	createFavoriteTable = `
        CREATE TABLE IF NOT EXISTS favorites (
            user_id UUID REFERENCES local_users(id) ON DELETE CASCADE,
            product_id UUID REFERENCES products(id) ON DELETE CASCADE,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
            PRIMARY KEY (user_id, product_id)
        )
    `

	createIndex = `
        CREATE INDEX ON products USING hnsw (embedding vector_cosine_ops)
        `

	createCleanupProductFunction = `
            
            CREATE OR REPLACE FUNCTION fn_cleanup_product_on_soft_delete()
            RETURNS TRIGGER AS $$
            BEGIN
                IF (NEW.status = 'deleted' AND OLD.status != 'deleted') THEN
                    DELETE FROM favorites WHERE product_id = NEW.id;
                    DELETE FROM user_product_interactions WHERE product_id = NEW.id;
                    UPDATE product_images 
                        SET deleted_at = NOW()
                        WHERE product_id = NEW.id AND deleted_at IS NULL;
                    RAISE NOTICE 'Product % soft-deleted, relations cleaned up.', NEW.id;
                END IF;
                RETURN NEW;
            END;
            $$ LANGUAGE plpgsql;

            
            DROP TRIGGER IF EXISTS trg_after_product_soft_delete ON products;

          
            CREATE TRIGGER trg_after_product_soft_delete
            AFTER UPDATE ON products
            FOR EACH ROW
            EXECUTE FUNCTION fn_cleanup_product_on_soft_delete();
        `
)
