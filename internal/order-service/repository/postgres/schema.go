package postgres

const (
	createOrdersTable = `
        CREATE TABLE IF NOT EXISTS orders (
            id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
            user_id UUID NOT NULL,
            total_price DECIMAL(12, 2) NOT NULL,
            -- Status: 1:Pending, 2:Paid, 3:Shipped, 4:Cancelled, 5:Completed
            status INT NOT NULL DEFAULT 1,
          
            shipping_address TEXT,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )  
    `

	createOrderItemsTable = `
		CREATE TABLE IF NOT EXISTS order_items (
            id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
            order_id UUID NOT NULL,
            product_id UUID NOT NULL,
            seller_id UUID NOT NULL, 
            quantity INT NOT NULL,
            product_name VARCHAR(255) NOT NULL, 
            product_image_url TEXT,             
            
            unit_price DECIMAL(12, 2) NOT NULL,
            status INT NOT NULL DEFAULT 1,
            
            CONSTRAINT fk_order
                FOREIGN KEY(order_id) 
                REFERENCES orders(id)
                ON DELETE CASCADE
        )  
	`

	createIndex = `
		CREATE INDEX idx_orders_user_id ON orders(user_id);
		CREATE INDEX idx_order_items_order_id ON order_items(order_id);
	`
)
