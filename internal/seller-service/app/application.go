// internal/user-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/seller-service/config"
	"marketplace/internal/seller-service/domain"
	"marketplace/internal/seller-service/infrastructure"
	"marketplace/internal/seller-service/pkg/graceful"
	"marketplace/internal/seller-service/repository/postgres"
	"marketplace/internal/seller-service/server"
	httptransport "marketplace/internal/seller-service/transport/http"
	"marketplace/pkg/messaging"
	pb "marketplace/pkg/proto/events"

	"time"
)

type App struct {
	cfg           config.Config
	server        *server.Server
	repository    domain.SellerRepository
	messaging     domain.Messaging
	cloudinarySvc domain.ImageService
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
		messaging:  container.messaging,
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
	repo          domain.SellerRepository
	server        *server.Server
	messaging     domain.Messaging
	cloudinarySvc domain.ImageService
}

func createMessagingConfig(cfg config.MessagingConfig) messaging.KafkaConfig {
	broker := cfg.Brokers[0]
	if broker == "" {
		broker = "localhost:29092"
	}
	kafkaBrokers := []string{broker}
	return messaging.KafkaConfig{
		Brokers:              kafkaBrokers,
		Topic:                "main-events", // Ana olay topic'i
		RetryTopic:           "main-events-retry",
		DLQTopic:             "main-events-dlq",
		ServiceType:          pb.ServiceType_SELLER_SERVICE,
		EnableRetry:          true,
		MaxRetries:           3,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []pb.MessageType{pb.MessageType_SELLER_APPROVED},
	}
}

func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	messagingConfig := createMessagingConfig(cfg.Messaging)

	messaging, err := messaging.NewKafkaClient(messagingConfig)
	if err != nil {
		return nil, fmt.Errorf("init kafka messaging: %w", err)
	}
	cloudinarySvc, err := infrastructure.NewCloudinaryService(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		return nil, fmt.Errorf("init cloudinary service: %w", err)
	}
	httpHandlers := httptransport.NewHandlers(repo, messaging,cloudinarySvc)
	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httpServer := server.New(serverCfg, router)

	return &container{
		repo:      repo,
		server:    httpServer,
		messaging: messaging,
	}, nil
}
