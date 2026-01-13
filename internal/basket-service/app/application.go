// internal/basket-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/basket-service/config"
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/pkg/graceful"
	"marketplace/internal/basket-service/repository/basket"
	"marketplace/internal/basket-service/repository/postgres"
	"marketplace/internal/basket-service/server"
	httptransport "marketplace/internal/basket-service/transport/http"

	"time"
)

type App struct {
	cfg                      config.Config
	server                   *server.Server
	BasketPostgresRepository domain.BasketPostgresRepository
	BasketRedisRepository    domain.BasketRedisRepository
}

func NewApp(cfg config.Config) (*App, error) {
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return &App{
		cfg:                      cfg,
		server:                   container.server,
		BasketPostgresRepository: container.BasketPostgresRepository,
		BasketRedisRepository:    container.BasketRedisRepository,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)

	log.Printf("starting user-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	if err := a.BasketRedisRepository.Close(); err != nil {
		return fmt.Errorf("failed to close redis repository: %w", err)
	}
	return a.BasketPostgresRepository.Close()
}

type container struct {
	BasketPostgresRepository domain.BasketPostgresRepository
	BasketRedisRepository    domain.BasketRedisRepository
	server                   *server.Server
}

func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	redisRepo, err := basket.NewBasketRedisRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init redis repository: %w", err)
	}

	httpHandlers := httptransport.NewHandlers(repo, redisRepo)
	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httpServer := server.New(serverCfg, router)

	return &container{
		BasketPostgresRepository: repo,
		BasketRedisRepository:    redisRepo,
		server:                   httpServer,
	}, nil
}
