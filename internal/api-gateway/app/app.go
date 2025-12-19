package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"marketplace/internal/api-gateway/cache"
	"marketplace/internal/api-gateway/config"
	"marketplace/internal/api-gateway/server"
	"marketplace/internal/user-service/pkg/graceful"
)

type App struct {
	cfg config.Config
	// Fiber        *fiber.App

	server *server.Server
	// Registry     *service.ServiceRegistry
	// RateLimiter  *limiter.RateLimiter
	// Metrics      *metrics.Metrics
	cacheManager *cache.CacheManager
}

func New(cfg config.Config) *App {

	// Wait, registry.go define NewServiceRegistry properly?
	// YES: func NewServiceRegistry() *ServiceRegistry
	cacheManager, err := cache.NewCacheManager(
		cfg.RedisCache.Addr,
		cfg.RedisCache.Password,
		cfg.RedisCache.DB,
		time.Duration(cfg.RedisCache.CacheTTL)*time.Second,
	)
	if err != nil {
		log.Fatalf("❌ Redis Cache başlatılamadı: %v", err)
	}
	log.Println("✅ Redis Cache başarıyla bağlandı")

	// Uygulama kapanırken cache'i kapat

	server := server.New(cfg, cacheManager)
	return &App{
		//Fiber:        f,
		// Registry:     registry,
		// RateLimiter:  rateLimiter,
		// Metrics:      metrics,
		server:       server,
		cacheManager: cacheManager,
	}
}
func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)

	log.Printf("starting user-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.cacheManager.Close()
}

func (a *App) RegisterService(name string, baseURLs []string, prefix string) error {
	return a.server.Registry.Register(name, baseURLs, prefix)
}
