// internal/paybent-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/payment-service/config"
	"marketplace/internal/payment-service/pkg/graceful"
	"marketplace/internal/payment-service/server"
	httptransport "marketplace/internal/payment-service/transport/http"
	"time"
)

type App struct {
	cfg    config.Config
	server *server.Server
}

func NewApp(cfg config.Config) (*App, error) {
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return &App{
		cfg:    cfg,
		server: container.server,
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
	return nil
}

type container struct {
	server *server.Server
}

func buildContainer(cfg config.Config) (*container, error) {

	httpHandlers := httptransport.NewHandlers()
	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httpServer := server.New(serverCfg, router)

	return &container{

		server: httpServer,
	}, nil
}
