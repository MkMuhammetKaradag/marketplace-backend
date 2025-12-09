package app

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"marketplace/internal/api-gateway/config"
	"marketplace/internal/api-gateway/handlers"
	"marketplace/internal/api-gateway/limiter"
	"marketplace/internal/api-gateway/metrics"
	"marketplace/internal/api-gateway/middleware"
	"marketplace/internal/api-gateway/service"
)

type App struct {
	Fiber       *fiber.App
	Registry    *service.ServiceRegistry
	RateLimiter *limiter.RateLimiter
	Metrics     *metrics.Metrics
}

func New() *App {
	// Initialize Components
	registry := service.NewServiceRegistry()
	// Wait, registry.go define NewServiceRegistry properly?
	// YES: func NewServiceRegistry() *ServiceRegistry

	registry = service.NewServiceRegistry()
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
	f.Use(middleware.AuthMiddleware(config.GetProtectedRoutes()))
	// 3. Rate Limit
	f.Use(middleware.RateLimitMiddleware(rateLimiter, metrics, config.GetDefaultRouteConfigs()))

	// Handlers
	proxyHandler := handlers.NewProxyHandler(registry, metrics)
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

	// Start Background Tasks
	registry.StartHealthChecks(15 * time.Second)
	rateLimiter.StartCleanup(5*time.Minute, 15*time.Minute)

	return &App{
		Fiber:       f,
		Registry:    registry,
		RateLimiter: rateLimiter,
		Metrics:     metrics,
	}
}

func (a *App) Run(addr string) error {
	return a.Fiber.Listen(addr)
}

func (a *App) RegisterService(name string, baseURLs []string, prefix string) error {
	return a.Registry.Register(name, baseURLs, prefix)
}
