// internal/user-service/app/app.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/user-service/config"
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/pkg/graceful"
	"marketplace/internal/user-service/repository/postgres"
	"marketplace/internal/user-service/server"
	httptransport "marketplace/internal/user-service/transport/http"
	"time"
)

type App struct {
	cfg        config.Config
	server     *server.Server
	repository domain.UserRepository
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
	repo   domain.UserRepository
	server *server.Server
}

func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	userService := domain.NewUserService(repo)
	httpHandlers := httptransport.NewHandlers(userService, repo)
	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httpServer := server.New(serverCfg, router)

	return &container{
		repo:   repo,
		server: httpServer,
	}, nil
}
