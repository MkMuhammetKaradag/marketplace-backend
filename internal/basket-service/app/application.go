// internal/basket-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/basket-service/config"
	"marketplace/internal/basket-service/domain"
	"marketplace/internal/basket-service/grpc_client"
	"marketplace/internal/basket-service/pkg/graceful"
	"marketplace/internal/basket-service/repository/basket"
	"marketplace/internal/basket-service/repository/postgres"
	"marketplace/internal/basket-service/server"
	grpctransport "marketplace/internal/basket-service/transport/grpc"
	httptransport "marketplace/internal/basket-service/transport/http"
	"marketplace/internal/basket-service/transport/kafka"
	messaginghandler "marketplace/internal/basket-service/transport/messaging"
	"time"
)

type App struct {
	cfg                      config.Config
	server                   *server.Server
	BasketPostgresRepository domain.BasketPostgresRepository
	BasketRedisRepository    domain.BasketRedisRepository
	consumer                 *kafka.Consumer
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
		consumer:                 container.consumer,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)
	a.consumer.Start(ctx)
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
	consumer                 *kafka.Consumer
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
	grpcAddress := fmt.Sprintf("localhost:%s", cfg.Server.GrpcProductPort)
	grpcProductClient, err := grpc_client.NewProductClient(grpcAddress)
	if err != nil {
		log.Fatalf("failed to initialise gRPC client: %v", err)
	}
	httpHandlers := httptransport.NewHandlers(repo, redisRepo, grpcProductClient)
	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		GrpcPort:     cfg.Server.GrpcPort,
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	messsagingHnadlers := messaginghandler.SetupMessageHandlers(redisRepo)
	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, messsagingHnadlers)
	if err != nil {
		return nil, fmt.Errorf("init kafka consumer: %w", err)
	}

	grpcHandler := grpctransport.NewBasketGrpcHandler(redisRepo)
	httpServer := server.New(serverCfg, router, grpcProductClient, grpcHandler)

	return &container{
		BasketPostgresRepository: repo,
		BasketRedisRepository:    redisRepo,
		server:                   httpServer,
		consumer:                 kafkaConsumer,
	}, nil
}
