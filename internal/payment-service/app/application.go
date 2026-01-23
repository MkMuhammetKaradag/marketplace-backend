// internal/paybent-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/payment-service/config"
	"marketplace/internal/payment-service/infrastructure/payment"
	"marketplace/internal/payment-service/pkg/graceful"
	"marketplace/internal/payment-service/repository/postgres"
	"marketplace/internal/payment-service/server"
	grpctransport "marketplace/internal/payment-service/transport/grpc"
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
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}
	stripeService := payment.NewStripeService(cfg.Stripe.SecretKey, cfg.Stripe.WebhookSecret)

	httpHandlers := httptransport.NewHandlers(stripeService)
	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port: cfg.Server.Port,

		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		GrpcPort:     cfg.Server.GrpcPort,
	}

	grpcHandler := grpctransport.NewProductGrpcHandler(repo, stripeService)

	httpServer := server.New(serverCfg, router, grpcHandler)

	return &container{

		server: httpServer,
	}, nil
}
