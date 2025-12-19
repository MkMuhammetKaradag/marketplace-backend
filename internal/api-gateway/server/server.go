package server

import (
	"fmt"
	"log"
	"marketplace/internal/api-gateway/cache"
	"marketplace/internal/api-gateway/config"
	"marketplace/internal/api-gateway/grpc_client"
	"marketplace/internal/api-gateway/handlers"
	"marketplace/internal/api-gateway/limiter"
	"marketplace/internal/api-gateway/metrics"
	"marketplace/internal/api-gateway/middleware"
	"marketplace/internal/api-gateway/service"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type RouteRegistrar interface {
	Register(app *fiber.App)
}

type Server struct {
	app         *fiber.App
	cfg         config.Config
	Registry    *service.ServiceRegistry
	RateLimiter *limiter.RateLimiter
	Metrics     *metrics.Metrics
}

func New(cfg config.Config, cacheManager *cache.CacheManager) *Server {
	registry := service.NewServiceRegistry()
	rateLimiter := limiter.NewRateLimiter()
	metrics := metrics.NewMetrics()
	// Initialize Fiber App
	f := fiber.New(fiber.Config{
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Global Middleware
	f.Use(recover.New())
	f.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} - ${latency}\n",
	}))
	f.Use(cors.New())

	// Custom Middleware
	// 1. Logging (Done by fiber/logger above roughly, or we can use custom if we want extraction)
	// 2. Auth
	f.Use(middleware.AuthMiddleware(config.GetProtectedRoutes(), cacheManager))
	// 3. Rate Limit
	f.Use(middleware.RateLimitMiddleware(rateLimiter, metrics, config.GetDefaultRouteConfigs()))

	// Handlers
	proxyHandler := handlers.NewProxyHandler(registry, metrics, cacheManager)
	manageHandler := handlers.NewManageHandler(registry, metrics)

	// Management Routes
	f.Get("/health", manageHandler.HealthCheck)
	f.Get("/metrics", manageHandler.GetMetrics)
	f.Get("/services", manageHandler.ListServices)
	f.Get("/simulate/login", manageHandler.SimulateLogin)

	// Proxy Route (Catch-all)
	// We matched "/" in original but `http.NewServeMux` matches prefixes.
	// In Fiber `*` works as wildcard.
	f.All("/*", proxyHandler.Handle)
	registry.StartHealthChecks(15 * time.Second)
	rateLimiter.StartCleanup(5*time.Minute, 15*time.Minute)
	return &Server{
		app:         f,
		cfg:         cfg,
		Registry:    registry,
		RateLimiter: rateLimiter,
		Metrics:     metrics,
	}
}

func (s *Server) Start() error {
	// 1. gRPC sunucusunu bir goroutine içinde başlatın
	// Fiber'in Listen() çağrısı bloklayıcı olduğu için bunu yapmalıyız.
	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gRPC sunucusu hatası: %v", err)
		}
	}()
	// 2. HTTP Fiber sunucusunu başlatın (Bu çağrı bloklayıcıdır)
	log.Printf("🌐 HTTP sunucusu %s adresinde dinliyor...", s.cfg.Server.Port)
	return s.app.Listen(s.Address())
}

func (s *Server) Shutdown(timeout time.Duration) error {

	return s.app.ShutdownWithTimeout(timeout)
}

func (s *Server) FiberApp() *fiber.App {
	return s.app
}

func (s *Server) Address() string {
	return fmt.Sprintf("0.0.0.0:%s", s.cfg.Server.Port)
}

func (s *Server) Run() error {
	grpcAddress := "localhost:3001" // Docker'da ise servis adı, yerelde ise localhost:50051

	if err := grpc_client.InitAuthClient(grpcAddress); err != nil {
		log.Fatalf("gRPC istemcisi başlatılamadı: %v", err)
		return err
	}
	return nil
}
