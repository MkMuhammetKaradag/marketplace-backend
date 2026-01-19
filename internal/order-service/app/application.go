// internal/order-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/order-service/config"
	"marketplace/internal/order-service/domain"
	"marketplace/internal/order-service/pkg/graceful"
	"marketplace/internal/order-service/repository/postgres"
	"marketplace/internal/order-service/server"
	httptransport "marketplace/internal/order-service/transport/http"
	"time"
)

type App struct {
	cfg        config.Config
	server     *server.Server
	repository domain.OrderRepository
}

func NewApp(cfg config.Config) (*App, error) {
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return &App{
		cfg:        cfg,
		server:     container.server,
		repository: container.repository,
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
	return a.repository.Close()
}

type container struct {
	server     *server.Server
	repository domain.OrderRepository
}

func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}
	httpHandlers := httptransport.NewHandlers(repo)
	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httpServer := server.New(serverCfg, router)

	return &container{

		server:     httpServer,
		repository: repo,
	}, nil
}
