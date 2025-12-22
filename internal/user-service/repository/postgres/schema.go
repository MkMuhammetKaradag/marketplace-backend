package postgres

const (
	// --- Genel Kullanıcı İzinleri (0-9) ---
	PermissionViewProduct   int64 = 1 << 0 // Ürünleri görüntüleyebilme
	PermissionWriteReview   int64 = 1 << 1 // Yorum ve değerlendirme yapabilme
	PermissionContactSeller int64 = 1 << 2 // Satıcıya mesaj atabilme
	PermissionPlaceOrder    int64 = 1 << 3 // Sipariş verebilme (Satın alma)

	// --- Satıcı İzinleri (10-19) ---
	PermissionManageOwnStore int64 = 1 << 10 // Kendi mağaza bilgilerini düzenleme
	PermissionAddProduct     int64 = 1 << 11 // Yeni ürün ekleyebilme
	PermissionEditProduct    int64 = 1 << 12 // Mevcut ürünlerini güncelleyebilme
	PermissionDeleteProduct  int64 = 1 << 13 // Ürünlerini silebilme/arşive alma
	PermissionManageOrders   int64 = 1 << 14 // Gelen siparişleri onaylama/kargolama
	PermissionViewAnalytics  int64 = 1 << 15 // Satış istatistiklerini görme

	// --- Moderasyon ve Destek İzinleri (20-29) ---
	PermissionApproveProducts int64 = 1 << 20 // Satıcıların ürünlerini yayına almadan önce onaylama
	PermissionManageDisputes  int64 = 1 << 21 // Alıcı-Satıcı arasındaki itirazları yönetme
	PermissionBanUsers        int64 = 1 << 22 // Kural ihlali yapanları yasaklama
	PermissionViewAllOrders   int64 = 1 << 23 // Sistemdeki tüm sipariş detaylarını görme

	// --- Finans ve Üst Yönetim İzinleri (30-39) ---
	PermissionManagePayments int64 = 1 << 30 // Ödeme geri iadeleri ve hakedişleri yönetme
	PermissionSetCommissions int64 = 1 << 31 // Kategori bazlı komisyon oranlarını belirleme
	PermissionManageRoles    int64 = 1 << 32 // Yeni roller ve yetkiler tanımlama
	PermissionAdministrator  int64 = 1 << 62 // TAM YETKİ (Sistem sahibi)
)
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
