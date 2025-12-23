package domain // İzin listesini bir dizi olarak tanımlıyoruz

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
	PermissionApproveProducts       int64 = 1 << 20 // Satıcıların ürünlerini yayına almadan önce onaylama
	PermissionManageDisputes        int64 = 1 << 21 // Alıcı-Satıcı arasındaki itirazları yönetme
	PermissionBanUsers              int64 = 1 << 22 // Kural ihlali yapanları yasaklama
	PermissionViewAllOrders         int64 = 1 << 23 // Sistemdeki tüm sipariş detaylarını görme
	PermissionApproveOrRejectSeller int64 = 1 << 24 // Satıcı onayını yönetme (onaylama veya reddetme)
	PermissionRemoveSeller          int64 = 1 << 25 // Satıcı onayını kaldırma
	PermissionManageReports         int64 = 1 << 26 // Raporları yönetme

	// --- Finans ve Üst Yönetim İzinleri (30-39) ---
	PermissionManagePayments int64 = 1 << 30 // Ödeme geri iadeleri ve hakedişleri yönetme
	PermissionSetCommissions int64 = 1 << 31 // Kategori bazlı komisyon oranlarını belirleme
	PermissionManageRoles    int64 = 1 << 32 // Yeni roller ve yetkiler tanımlama
	PermissionAdministrator  int64 = 1 << 62 // TAM YETKİ (Sistem sahibi)
)

var AllPermissionsList = []int64{
	PermissionViewProduct,
	PermissionWriteReview,
	PermissionContactSeller,
	PermissionPlaceOrder,
	PermissionManageOwnStore,
	PermissionAddProduct,
	PermissionEditProduct,
	PermissionDeleteProduct,
	PermissionManageOrders,
	PermissionViewAnalytics,
	PermissionApproveProducts,
	PermissionManageDisputes,
	PermissionBanUsers,
	PermissionViewAllOrders,
	PermissionApproveOrRejectSeller,
	PermissionRemoveSeller,
	PermissionManageReports,
	PermissionManagePayments,
	PermissionSetCommissions,
	PermissionManageRoles,
	PermissionAdministrator,
}


var ValidPermissionsMask int64

func init() {
	for _, p := range AllPermissionsList {
		ValidPermissionsMask |= p
	}
}


func IsValidPermission(p int64) bool {
	return (p & ^ValidPermissionsMask) == 0
}
