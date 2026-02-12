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

	// Kafka'yı başlat
	a.consumer.Start(ctx)

	// Logu düzelttik
	log.Printf("starting basket-service on %s (gRPC: %s)", a.cfg.Server.Port, a.cfg.Server.GrpcPort)

	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repositories")
	// Repoları güvenli kapatma
	_ = a.BasketRedisRepository.Close()
	return a.BasketPostgresRepository.Close()
}

type container struct {
	BasketPostgresRepository domain.BasketPostgresRepository
	BasketRedisRepository    domain.BasketRedisRepository
	server                   *server.Server
	consumer                 *kafka.Consumer
}

func buildContainer(cfg config.Config) (*container, error) {
	// 1. Repositories
	postgresRepo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	redisRepo, err := basket.NewBasketRedisRepository(cfg)
	if err != nil {
		return nil, err
	}

	// 2. gRPC Client (Product Service'e bağlanmak için)
	grpcAddress := fmt.Sprintf("localhost:%s", cfg.Server.GrpcProductPort)
	productClient, err := grpc_client.NewProductClient(grpcAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to init gRPC product client: %w", err)
	}

	// 3. Handlers & Router (Grup yapısına hazırlık)
	// Diğer servislerdeki gibi NewHandlers içinde UseCase'leri kuracağız
	h := httptransport.NewHandlers(postgresRepo, redisRepo, productClient)
	router := httptransport.NewRouter(h)

	// 4. Messaging
	msgHandlers := messaginghandler.SetupMessageHandlers(redisRepo)
	consumer, err := kafka.NewConsumer(cfg.Messaging, msgHandlers)
	if err != nil {
		return nil, err
	}

	// 5. Server
	grpcHandler := grpctransport.NewBasketGrpcHandler(redisRepo)
	s := server.New(getServerConfig(cfg), router, grpcHandler)

	return &container{
		BasketPostgresRepository: postgresRepo,
		BasketRedisRepository:    redisRepo,
		server:                   s,
		consumer:                 consumer,
	}, nil
}
func getServerConfig(cfg config.Config) server.Config {
	return server.Config{
		Port:         cfg.Server.Port,
		GrpcPort:     cfg.Server.GrpcPort,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
