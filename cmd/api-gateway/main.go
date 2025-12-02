package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/google/uuid" // JWT yerine basit oturum ID'leri kullandığımız için
	"golang.org/x/time/rate"
)

// --- YAPILANDIRMA (Configuration) ---
// Hassas verileri çevre değişkenlerinden okumak en iyisidir.
const (
	// Gateway'in iç servislerle iletişim kurarken kullandığı gizli anahtar
	InternalGatewayHeader = "X-API-Key"
	// os.Getenv ile çevresel değişkenden alınmalı. Varsayılan (default) değer.
	InternalGatewaySecret = "GATEWAY_SECRET_KEY"

	SessionCookieName = "session_id"
	DefaultTimeout    = 30 * time.Second
	GatewayPort       = ":8080"
)

// --- SERVİS KAYIT SİSTEMİ ve YÜK DENGELEME (Service Registry & Load Balancing) ---

// ServiceHealth: Servisin anlık sağlık durumunu tutar.
type ServiceHealth struct {
	Healthy   bool
	LastCheck time.Time
	FailCount int32
	mu        sync.RWMutex
}

// Service: Kayıtlı bir arka uç servisini temsil eder.
type Service struct {
	Name       string
	BaseURLs   []string // Yük dengeleme için birden fazla URL
	PathPrefix string
	Health     *ServiceHealth
	Timeout    time.Duration
	// Load Balancing için kullanılacak indeks
	nextIndex uint64
}

// ServiceRegistry: Servislerin kaydını ve yönetimini sağlar.
type ServiceRegistry struct {
	services sync.Map // map[string]*Service
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{}
}

// Register: Bir servisi birden fazla adresiyle kaydeder.
func (sr *ServiceRegistry) Register(name string, baseURLs []string, pathPrefix string) error {
	if len(baseURLs) == 0 {
		return fmt.Errorf("en az bir BaseURL gerekli")
	}

	service := &Service{
		Name:       name,
		BaseURLs:   baseURLs,
		PathPrefix: pathPrefix,
		Health: &ServiceHealth{
			Healthy:   true,
			LastCheck: time.Now(),
		},
		Timeout:   DefaultTimeout,
		nextIndex: 0,
	}
	sr.services.Store(name, service)
	log.Printf("✅ Servis kaydedildi: %s -> %v (prefix: %s)", name, baseURLs, pathPrefix)
	return nil
}

// GetByPath: Gelen isteğin yoluna göre ilgili servisi bulur. En uzun eşleşme önceliklidir.
func (sr *ServiceRegistry) GetByPath(path string) (*Service, bool) {
	var found *Service
	longestPrefixLen := 0

	sr.services.Range(func(key, value interface{}) bool {
		service := value.(*Service)
		// Yolu kontrol et ve en uzun (en spesifik) prefix'i bul
		if strings.HasPrefix(path, service.PathPrefix) {
			if len(service.PathPrefix) > longestPrefixLen {
				found = service
				longestPrefixLen = len(service.PathPrefix)
			}
		}
		return true
	})
	return found, found != nil
}
func (sr *ServiceRegistry) List() []*Service {
	var services []*Service
	// sync.Map.Range, harita üzerinde döngü kurmanın eş zamanlı güvenli yoludur.
	sr.services.Range(func(key, value interface{}) bool {
		// Her değeri *Service türüne dönüştürüp listeye ekle
		services = append(services, value.(*Service))
		return true // Döngüye devam et
	})
	return services
}

// GetNextBaseURL: Servisin sağlıklı durumdaki bir sonraki URL'sini Round-Robin ile döndürür.
func (s *Service) GetNextBaseURL() (string, bool) {
	s.Health.mu.RLock()
	// Servis sağlıklı değilse yük dengeleme yapmaya gerek yok
	if !s.Health.Healthy {
		s.Health.mu.RUnlock()
		return "", false
	}
	s.Health.mu.RUnlock()

	// Atomik olarak bir sonraki indeksi al ve artır (eş zamanlı güvenli sayım)
	index := atomic.AddUint64(&s.nextIndex, 1) - 1
	// Modulo ile BaseURLs dizisinin sınırları içinde kal
	urlIndex := index % uint64(len(s.BaseURLs))

	return s.BaseURLs[urlIndex], true
}

// StartHealthChecks: Periyodik sağlık kontrol mekanizmasını başlatır.
func (sr *ServiceRegistry) StartHealthChecks(interval time.Duration) {
	// ... (Mevcut health check mantığı)
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			sr.services.Range(func(key, value interface{}) bool {
				service := value.(*Service)
				go sr.checkHealth(service)
				return true
			})
		}
	}()
	log.Printf("🏥 Health check başlatıldı (interval: %v)", interval)
}

// checkHealth: Servisin her bir örneğini kontrol eder (Circuit Breaker mantığı içerir)
func (sr *ServiceRegistry) checkHealth(service *Service) {
	// Birden fazla URL'den birinin sağlıklı olması yeterli olabilir
	allHealthy := true

	service.Health.mu.Lock()
	defer service.Health.mu.Unlock()

	service.Health.LastCheck = time.Now()

	client := &http.Client{Timeout: 5 * time.Second}

	// Servisin tüm örneklerini kontrol et
	for _, baseURL := range service.BaseURLs {
		resp, err := client.Get(baseURL + "/health")

		if err != nil || resp.StatusCode != http.StatusOK {
			// Başarısızlık durumunda Circuit Breaker mantığı:
			// FailCount'u atomik olarak artır
			atomic.AddInt32(&service.Health.FailCount, 1)
			allHealthy = false
		} else {
			if resp != nil {
				resp.Body.Close()
			}
			// Başarılı olursa FailCount'u sıfırla
			atomic.StoreInt32(&service.Health.FailCount, 0)
		}
	}

	// Devre Kesici Kontrolü
	if !allHealthy && atomic.LoadInt32(&service.Health.FailCount) >= 3 {
		if service.Health.Healthy {
			log.Printf("🔴 Servis DOWN (Circuit Open): %s", service.Name)
			service.Health.Healthy = false // Devreyi aç
		}
	} else if allHealthy {
		if !service.Health.Healthy {
			log.Printf("🟢 Servis UP (Circuit Closed): %s", service.Name)
		}
		service.Health.Healthy = true
	}
}

// --- RATE LIMITING (Hız Sınırlama) ---

// RouteConfig: Bir yol veya prefix için belirlenmiş hız limitleri
type RouteConfig struct {
	GlobalLimit rate.Limit // Tüm kullanıcılar için toplam limit (request/s)
	GlobalBurst int        // Patlama isteği sayısı
	UserLimit   rate.Limit // Oturum/IP başına limit (request/s)
	UserBurst   int
}

// RateLimiter ve LimiterEntry yapıları önceki kodunuzdan aynen alınmıştır,
// eş zamanlı güvenli ve temizleme mekanizmalı olduğu için.

type LimiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

type RateLimiter struct {
	limiters sync.Map // map[string]*LimiterEntry
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{}
}

// GetLimiter: Belirtilen anahtar için bir rate.Limiter döndürür veya oluşturur.
func (rl *RateLimiter) GetLimiter(key string, r rate.Limit, b int) *rate.Limiter {
	// ... (Mevcut GetLimiter mantığı)
	now := time.Now()
	if entry, ok := rl.limiters.Load(key); ok {
		limiterEntry := entry.(*LimiterEntry)
		limiterEntry.lastAccess = now
		return limiterEntry.limiter
	}
	newLimiter := rate.NewLimiter(r, b)
	entry := &LimiterEntry{
		limiter:    newLimiter,
		lastAccess: now,
	}
	actual, loaded := rl.limiters.LoadOrStore(key, entry)
	if loaded {
		return actual.(*LimiterEntry).limiter
	}
	return newLimiter
}

// StartCleanup: Kullanılmayan limiteleri periyodik olarak siler.
func (rl *RateLimiter) StartCleanup(interval, maxAge time.Duration) {
	// ... (Mevcut StartCleanup mantığı)
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			count := 0
			now := time.Now()
			rl.limiters.Range(func(key, value interface{}) bool {
				entry := value.(*LimiterEntry)
				if now.Sub(entry.lastAccess) > maxAge {
					rl.limiters.Delete(key)
					count++
				}
				return true
			})
			if count > 0 {
				log.Printf("🧹 Temizlik: %d limiter silindi", count)
			}
		}
	}()
}

// --- METRİKLER (Metrics) ---
// Metrik yapısı önceki kodunuzdan aynen alınmıştır.

type Metrics struct {
	mu                  sync.RWMutex
	totalRequests       int64
	successRequests     int64
	failedRequests      int64
	rateLimitedRequests map[string]int64
	requestsByService   map[string]int64
	requestsByPath      map[string]int64
	lastReset           time.Time
}

func NewMetrics() *Metrics {
	// ... (Metrik başlangıç mantığı)
	return &Metrics{
		rateLimitedRequests: make(map[string]int64),
		requestsByService:   make(map[string]int64),
		requestsByPath:      make(map[string]int64),
		lastReset:           time.Now(),
	}
}

// IncrementTotal, IncrementSuccess, GetStats gibi metrik fonksiyonları önceki kodunuzdan aynen kullanılır.

func (m *Metrics) IncrementTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.totalRequests++
}

func (m *Metrics) IncrementSuccess() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.successRequests++
}

func (m *Metrics) IncrementFailed() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failedRequests++
}

func (m *Metrics) IncrementRateLimit(limitType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.rateLimitedRequests[limitType]++
}

func (m *Metrics) IncrementService(service string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestsByService[service]++
}

func (m *Metrics) IncrementPath(path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestsByPath[path]++
}

func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	uptime := time.Since(m.lastReset)
	rps := float64(m.totalRequests) / uptime.Seconds()

	successRate := 0.0
	if m.totalRequests > 0 {
		successRate = float64(m.successRequests) / float64(m.totalRequests) * 100
	}
	return map[string]interface{}{
		"total_requests":        m.totalRequests,
		"success_requests":      m.successRequests,
		"failed_requests":       m.failedRequests,
		"success_rate":          fmt.Sprintf("%.2f%%", successRate),
		"rate_limited_requests": m.rateLimitedRequests,
		"requests_by_service":   m.requestsByService,
		"requests_by_path":      m.requestsByPath,
		"uptime_seconds":        uptime.Seconds(),
		"requests_per_second":   rps,
	}
}

// --- GATEWAY ---

type Gateway struct {
	registry     *ServiceRegistry
	rateLimiter  *RateLimiter
	metrics      *Metrics
	routeConfigs map[string]RouteConfig // Rate limit kuralları
	protected    map[string]bool        // Koruma gerektiren yollar
}

func NewGateway() *Gateway {
	// 10.0/60 = 10 istek/dakika (≈ 0.16 rps)
	return &Gateway{
		registry:    NewServiceRegistry(),
		rateLimiter: NewRateLimiter(),
		metrics:     NewMetrics(),
		routeConfigs: map[string]RouteConfig{
			// /users servisi için genel kural (prefix eşleşmesi)
			"/users": {
				GlobalLimit: 10.0 / 60, GlobalBurst: 10, // Max 200 istek/dakika
				UserLimit: 3.0 / 60, UserBurst: 3, // Kullanıcı max 20 istek/dakika
			},
			// /test servisi için genel kural
			"/test": {
				GlobalLimit: 20.0 / 60, GlobalBurst: 5,
				UserLimit: 3.0 / 60, UserBurst: 2,
			},
			// /test/hello için özel kural (tam yol eşleşmesi, en yüksek öncelik)
			"/test/hello": {
				GlobalLimit: 2.0 / 60, GlobalBurst: 2,
				UserLimit: 1.0 / 60, UserBurst: 1, // Kullanıcı max 1 istek/dakika
			},
			"/chat": {
				GlobalLimit: 50.0 / 60, GlobalBurst: 10,
				UserLimit: 5.0 / 60, UserBurst: 3,
			},
			// Varsayılan kural (eşleşmeyen tüm yollar)
			"default": {
				GlobalLimit: 1.0 / 60, GlobalBurst: 1,
				UserLimit: 1.0 / 60, UserBurst: 1,
			},
		},
		// Çerez/JWT gerektiren yollar
		protected: map[string]bool{
			"/users/profile": true,
			"/users/list":    true,
			"/test/hello":    true,
		},
	}
}

// --- MIDDLEWARE'LER ---

// corsMiddleware: CORS başlıklarını ekler ve OPTIONS isteklerini sonlandırır.
func (g *Gateway) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		// "Authorization" başlığını da kabul et
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		w.Header().Set("Access-Control-Expose-Headers", "X-RateLimit-Limit, X-RateLimit-Remaining, Retry-After")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// loggingMiddleware: İstek başlangıç ve bitiş loglarını tutar.
func (g *Gateway) loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		clientID := extractClientIdentifier(r)
		log.Printf("→ %s %s [%s]", r.Method, r.URL.Path, clientID)

		next(w, r)

		log.Printf("← %s %s [%dms]", r.Method, r.URL.Path, time.Since(start).Milliseconds())
	}
}

// authMiddleware: Oturum çerezi veya Authorization başlığını kontrol eder.
func (g *Gateway) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Sadece korumalı yollar için kontrol yap
		if g.protected[r.URL.Path] {
			isAuthenticated := false

			// 1. Çerez Kontrolü
			if _, err := r.Cookie(SessionCookieName); err == nil {
				isAuthenticated = true
			}

			// 2. Authorization Header Kontrolü (Bearer Token veya JWT)
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") && len(authHeader) > 7 {
				// Gerçek bir uygulamada burada token doğrulaması yapılır.
				// Şimdilik sadece varlığını kontrol ediyoruz.
				isAuthenticated = true
			}

			if !isAuthenticated {
				log.Printf("🔒 Yetkisiz erişim: %s", r.URL.Path)
				respondJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "Kimlik doğrulama (Oturum/Token) gerekli",
				})
				return
			}
		}
		next(w, r)
	}
}

// getRouteConfig: En spesifik (en uzun prefix/tam yol) rate limit kuralını bulur.
func (g *Gateway) getRouteConfig(path string) RouteConfig {
	config := g.routeConfigs["default"]
	longestMatchLen := 0

	// 1. Tam yol eşleşmesi kontrolü (En yüksek öncelik)
	if c, exists := g.routeConfigs[path]; exists {
		return c
	}

	// 2. Prefix eşleşmesi kontrolü
	for route, c := range g.routeConfigs {
		// "default" kuralı prefix olarak sayılmaz
		if route != "default" && strings.HasPrefix(path, route) {
			// En uzun prefix eşleşmesini bul
			if len(route) > longestMatchLen {
				config = c
				longestMatchLen = len(route)
			}
		}
	}

	return config
}

// rateLimitMiddleware: Hız sınırlama kurallarını uygular.
func (g *Gateway) rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		clientID := extractClientIdentifier(r)
		g.metrics.IncrementTotal()
		g.metrics.IncrementPath(path)

		config := g.getRouteConfig(path)

		// 2. KULLANICI BAŞINA LİMİT (User-path limit) - YENİ ALLOW() KULLANIMI
		if config.UserLimit > 0 {
			limiter := g.rateLimiter.GetLimiter("user:"+clientID+":"+path, config.UserLimit, config.UserBurst)

			// Allow() anında kontrol eder ve tokenı tüketir (Reserve()'dan farklıdır).
			if !limiter.Allow() {
				g.metrics.IncrementRateLimit("user-path")
				log.Printf("⛔ Rate limit (User): %s -> %s", clientID, path)

				// X-RateLimit-Limit başlıklarını buraya taşıyın (Allow() kullanınca Reserve() yok)
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", config.UserLimit*60))
				w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Tokens()))

				// Not: Allow() kullanırken Retry-After hesaplamak zordur.
				respondJSON(w, http.StatusTooManyRequests, map[string]string{
					"error": "Çok fazla istek",
					"type":  "user-path",
				})
				return // 🛑 KESİNLİKLE DÖN!
			}
			// İstek başarılı olduysa, X-RateLimit başlıklarını burada ayarlayın
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%.0f", config.UserLimit*60))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Tokens()))
		}

		// 3. GLOBAL LİMİT (Global path limit) - Allow() kullanılıyordu, doğru.
		limiter := g.rateLimiter.GetLimiter("global:"+path, config.GlobalLimit, config.GlobalBurst)
		if !limiter.Allow() {
			g.metrics.IncrementRateLimit("global-path")
			log.Printf("⛔ Rate limit (Global): %s", path)
			respondJSON(w, http.StatusTooManyRequests, map[string]string{
				"error": "Sistem yoğunluğu",
				"type":  "global-path",
			})
			return // 🛑 KESİNLİKLE DÖN!
		}

		next(w, r)
	}
}

// --- PROXY ve YÜK DENGELEME ---

// proxyHandler: İsteği yönlendirir, yük dengeleme ve devre kesici kontrolü yapar.
func (g *Gateway) proxyHandler(w http.ResponseWriter, r *http.Request) {
	service, ok := g.registry.GetByPath(r.URL.Path)
	if !ok {
		respondJSON(w, http.StatusNotFound, map[string]string{
			"error": "Servis bulunamadı",
		})
		return
	}

	// Devre Kesici Kontrolü
	service.Health.mu.RLock()
	healthy := service.Health.Healthy
	service.Health.mu.RUnlock()

	if !healthy {
		g.metrics.IncrementFailed()
		log.Printf("❌ Circuit Open: %s", service.Name)
		respondJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error":   "Servis şu anda kullanılamıyor (Circuit Breaker Açık)",
			"service": service.Name,
		})
		return
	}

	// Yük Dengeleme (Load Balancing)
	targetBaseURL, ok := service.GetNextBaseURL()
	if !ok {
		// Bu aslında Circuit Breaker kontrolünden sonra nadiren olmalı
		g.metrics.IncrementFailed()
		respondJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error":   "Sağlıklı servis örneği bulunamadı",
			"service": service.Name,
		})
		return
	}

	g.metrics.IncrementService(service.Name)

	targetPath := strings.TrimPrefix(r.URL.Path, service.PathPrefix)
	targetURL := targetBaseURL + targetPath
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	if err := proxyRequest(targetURL, w, r, service.Timeout); err != nil {
		g.metrics.IncrementFailed()
		log.Printf("❌ Proxy error [%s]: %v", service.Name, err)
		respondJSON(w, http.StatusBadGateway, map[string]string{
			"error": "Arka uç servis hatası",
		})
		return
	}
	g.metrics.IncrementSuccess()
}

// proxyRequest: Hedef URL'ye isteği yönlendirir.
func proxyRequest(targetURL string, w http.ResponseWriter, r *http.Request, timeout time.Duration) error {
	// İstek gövdesini kopyala (io.ReadAll istek gövdesini tüketir)
	bodyBytes, _ := io.ReadAll(r.Body)
	// Yeni istek oluştur (gövdeyi tekrar okumak için bytes.NewBuffer kullan)
	proxyReq, err := http.NewRequest(r.Method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}

	// Başlıkları kopyala
	for key, values := range r.Header {
		if key != "Host" { // Host başlığı hedef URL'ye ayarlanmalı
			for _, value := range values {
				proxyReq.Header.Add(key, value)
			}
		}
	}

	// İç iletişim anahtarını ekle
	proxyReq.Header.Set(InternalGatewayHeader, InternalGatewaySecret)
	// Gerçek istemci IP'sini arka uca ilet
	proxyReq.Header.Set("X-Forwarded-For", r.RemoteAddr)

	// İstemci timeout ile isteği gönder
	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(proxyReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Yanıt başlıklarını ve durumu kopyala
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)

	// Yanıt gövdesini kopyala
	_, err = io.Copy(w, resp.Body)
	return err
}

// --- HELPER FONKSİYONLAR ---

// extractClientIdentifier: Rate Limit için benzersiz istemci ID'sini alır. (Çerez > IP)
func extractClientIdentifier(r *http.Request) string {
	// Oturum çerezi varsa, çerez ID'sini kullan
	if cookie, err := r.Cookie(SessionCookieName); err == nil && cookie.Value != "" {
		return "session:" + cookie.Value
	}
	// Authorization başlığı varsa, token'ın ilk 8 karakterini kullan
	if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") && len(authHeader) > 7 {
		// JWT'nin benzersiz bir parçasını kullanmak en iyisidir, burada bir UUID oluşturuyoruz.
		// Gerçek bir uygulamada JWT'den kullanıcı ID'si alınır.
		return "token:" + authHeader[7:15]
	}
	// IP adresi yoksa
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = strings.Split(r.RemoteAddr, ":")[0]
	}
	return "ip:" + ip
}

// respondJSON: HTTP yanıtını JSON formatında hazırlar.
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// Hata durumunda bile JSON döndürülmesini sağlar
	json.NewEncoder(w).Encode(data)
}

// --- HANDLERS (Yönetim Uç Noktaları) ---

// healthHandler: Gateway ve kayıtlı servislerin genel sağlık durumunu döndürür.
func (g *Gateway) healthHandler(w http.ResponseWriter, r *http.Request) {
	services := g.registry.List()
	serviceHealth := make(map[string]interface{})

	for _, svc := range services {
		svc.Health.mu.RLock()
		serviceHealth[svc.Name] = map[string]interface{}{
			"healthy":    svc.Health.Healthy,
			"last_check": svc.Health.LastCheck.Format(time.RFC3339),
			"fail_count": atomic.LoadInt32(&svc.Health.FailCount),
			"base_urls":  svc.BaseURLs,
		}
		svc.Health.mu.RUnlock()
	}
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"gateway":   "healthy",
		"services":  serviceHealth,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// metricsHandler: Performans metriklerini döndürür.
func (g *Gateway) metricsHandler(w http.ResponseWriter, r *http.Request) {
	stats := g.metrics.GetStats()
	respondJSON(w, http.StatusOK, stats)
}

// servicesHandler: Kayıtlı servislerin listesini ve temel bilgilerini döndürür.
func (g *Gateway) servicesHandler(w http.ResponseWriter, r *http.Request) {
	services := g.registry.List()
	serviceList := make([]map[string]interface{}, 0, len(services))

	for _, svc := range services {
		svc.Health.mu.RLock()
		serviceList = append(serviceList, map[string]interface{}{
			"name":        svc.Name,
			"base_urls":   svc.BaseURLs,
			"path_prefix": svc.PathPrefix,
			"healthy":     svc.Health.Healthy,
			"next_index":  atomic.LoadUint64(&svc.nextIndex), // Yük dengeleme indeksini göster
		})
		svc.Health.mu.RUnlock()
	}

	respondJSON(w, http.StatusOK, serviceList)
}

// simulateAuthHandler: Basit bir oturum çerezi oluşturarak kimlik doğrulama simülasyonu yapar.
func (g *Gateway) simulateAuthHandler(w http.ResponseWriter, r *http.Request) {
	// Yeni bir benzersiz oturum ID'si oluştur
	sessionID := uuid.New().String()

	// Çerezi ayarla (güvenli ayarlar: HttpOnly, Secure vb. eklenmelidir)
	http.SetCookie(w, &http.Cookie{
		Name:    SessionCookieName,
		Value:   sessionID,
		Expires: time.Now().Add(24 * time.Hour),
		Path:    "/",
		// Güvenlik için HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode ayarları önemlidir
	})

	respondJSON(w, http.StatusOK, map[string]string{
		"message":    "Oturum başarıyla oluşturuldu.",
		"session_id": sessionID,
		"warning":    "Bu sadece bir simülasyondur.",
	})
}

// --- MAIN ---

func main() {
	gateway := NewGateway()

	// Servisleri kaydet (Yük dengeleme simülasyonu için birden fazla adres)
	// Not: Bu adreslerde gerçek servislerin çalışıyor olması gerekir.
	gateway.registry.Register("user-service", []string{"http://localhost:8081", "http://localhost:8083"}, "/users")
	gateway.registry.Register("auth-service", []string{"http://localhost:8084"}, "/auth")
	gateway.registry.Register("test-service", []string{"http://localhost:8082"}, "/test")
	gateway.registry.Register("chat-service", []string{"http://localhost:8085"}, "/chat")

	// Health checks ve cleanup başlat
	gateway.registry.StartHealthChecks(15 * time.Second)            // Daha sık kontrol
	gateway.rateLimiter.StartCleanup(5*time.Minute, 15*time.Minute) // Limiteleri temizle

	// Router
	mux := http.NewServeMux()

	// Yönetim Endpoints
	mux.HandleFunc("/health", gateway.healthHandler)
	mux.HandleFunc("/metrics", gateway.metricsHandler)
	mux.HandleFunc("/services", gateway.servicesHandler)
	// Simülasyon Endpointi (Test amaçlı çerez oluşturmak için)
	mux.HandleFunc("/simulate/login", gateway.simulateAuthHandler)

	// Main handler chain
	handler := http.HandlerFunc(gateway.proxyHandler)
	// Middleware zincirini içeriden dışarıya doğru uygula (Önce Proxy, Sonra Limit, Sonra Auth...)
	handler = gateway.rateLimitMiddleware(handler) // 3. Hız Sınırlama
	handler = gateway.authMiddleware(handler)      // 2. Kimlik Doğrulama
	handler = gateway.loggingMiddleware(handler)   // 1. Loglama
	handler = gateway.corsMiddleware(handler)      // 0. CORS (En Dış Katman)

	// Eşleşmeyen tüm yolları proxy zincirine yönlendir
	mux.Handle("/", handler)

	// Server
	server := &http.Server{
		Addr:         GatewayPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Gateway başlatıldı: http://localhost%s", GatewayPort)
	log.Printf("ℹ️  Kullanım:")
	log.Printf("  - /users/profile -> user-service'e yönlendirilir (Kimlik doğrulama gereklidir)")
	log.Printf("  - /test/hello    -> test-service'e yönlendirilir (Çok sıkı Rate Limit)")
	log.Printf("  - /simulate/login -> Test amaçlı oturum çerezi oluşturur")

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	<-sigChan
	log.Println("\n🛑 Shutdown başlatılıyor...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Shutdown error: %v", err)
	}
	log.Println("✅ Gateway kapatıldı")
}
