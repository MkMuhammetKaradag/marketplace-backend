// internal/product-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/product-service/config"
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/repository/postgres"
	"marketplace/internal/product-service/server"
	httptransport "marketplace/internal/product-service/transport/http"
	"time"
)

type App struct {
	cfg        config.Config
	server     *server.Server
	repository domain.ProductRepository
}

func NewApp(cfg config.Config) (*App, error) {
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return &App{
		cfg:        cfg,
		server:     container.server,
		repository: container.repo,
	}, nil
}

func (a *App) Start() error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Printf("starting user-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.repository.Close()
}

type container struct {
	repo   domain.ProductRepository
	server *server.Server
}

func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	productService := domain.NewProductService(repo)
	httpHandlers := httptransport.NewHandlers(productService, repo)

	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		GrpcPort:     cfg.Server.GrpcPort,
	}

	httpServer := server.New(serverCfg, router)

	return &container{
		repo:   repo,
		server: httpServer,
	}, nil
}
