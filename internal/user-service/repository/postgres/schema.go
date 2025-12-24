package postgres

const (
	
	createUsersTable = `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    		username VARCHAR(50) NOT NULL UNIQUE,
			first_name VARCHAR(50),
			last_name VARCHAR(50),
			phone_number VARCHAR(15),
			avatar_url TEXT,
			cover_url TEXT,
			bio TEXT,
		email VARCHAR(100) NOT NULL UNIQUE,
			

    		
    		password TEXT NOT NULL,
    		is_active BOOLEAN DEFAULT false,
    		is_email_verified BOOLEAN DEFAULT false,
			activation_code VARCHAR(6),
			activation_id UUID  DEFAULT gen_random_uuid(),
			activation_expiry TIMESTAMP WITH TIME ZONE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			deleted_at TIMESTAMP WITH TIME ZONE,
			failed_login_attempts INT DEFAULT 0,
			account_locked BOOLEAN DEFAULT false,
			lock_until TIMESTAMP WITH TIME ZONE,
			last_login TIMESTAMP WITH TIME ZONE,
			CONSTRAINT check_activation_code CHECK (
				activation_code ~ '^[0-9]{6}$'
			)
		)`
	createUserAddressesTable = `
			CREATE TABLE IF NOT EXISTS user_addresses (
				id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
				user_id UUID REFERENCES users(id) ON DELETE CASCADE,
				title VARCHAR(50) NOT NULL,
				address_line TEXT NOT NULL,
				city VARCHAR(50) NOT NULL,
				country VARCHAR(50) NOT NULL,
				zip_code VARCHAR(10),
				is_default BOOLEAN DEFAULT false,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			)`
	createForgotPasswordsTable = `
		CREATE TABLE IF NOT EXISTS forgot_passwords (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID REFERENCES users(id) ON DELETE CASCADE,
			attempt_count INT DEFAULT 0,
			expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`

	createRolesTable = `
		CREATE TABLE IF NOT EXISTS roles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			created_by UUID NULL REFERENCES users(id),
			name VARCHAR(100) NOT NULL UNIQUE,
			color VARCHAR(7) DEFAULT '#B9BBBE',
			position INT DEFAULT 0, -- Hiyerarşi için (0 en alt)
			permissions BIGINT DEFAULT 0,
			is_mentionable BOOLEAN DEFAULT TRUE,
			is_hoisted BOOLEAN DEFAULT FALSE,
			is_managed BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()	
		)`
	createdUserRolesTable = `
		CREATE TABLE IF NOT EXISTS user_roles (
			user_id UUID NOT NULL REFERENCES users(id),
			role_id UUID NOT NULL REFERENCES roles(id),
			assigned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			PRIMARY KEY (user_id, role_id)
		)`

	createDefaultRoles = `
		INSERT INTO roles (name, color, created_by, is_mentionable, is_hoisted, permissions, is_managed)
		VALUES
			-- 1. BUYER (Sadece izleme, yorum ve sipariş: 1|2|4|8 = 15)
			('Buyer', '#3498DB', NULL, TRUE, FALSE, 15, TRUE),

			-- 2. SELLER (Buyer yetkileri + Mağaza ve Ürün yönetimi: 15|1024|2048|4096|8192|16384|32768 = 64527)
			('Seller', '#2ECC71', NULL, TRUE, TRUE, 64527, TRUE),

			-- 3. MODERATOR (İçerik onaylama ve uyuşmazlık çözme: 1|1048576|2097152|8388608 = 11534337)
			('Moderator', '#F1C40F', NULL, TRUE, TRUE, 11534337, TRUE),

			-- 4. ADMIN (Tüm yetkiler: 4611686018427387904)
			('Admin', '#E74C3C', NULL, TRUE, TRUE, 4611686018427387904, TRUE)ON CONFLICT (name) DO NOTHING
		`
)
