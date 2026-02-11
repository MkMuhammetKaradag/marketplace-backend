// internal/seller-service/app/application.go
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
		cfg:           cfg,
		server:        container.server,
		repository:    container.repo,
		messaging:     container.messaging,
		cloudinarySvc: container.cloudinarySvc,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)

	log.Printf("starting seller-service on %s", a.server.Address())
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

func getKafkaSettings(cfg config.MessagingConfig) messaging.KafkaConfig {
	broker := "localhost:29092"
	if len(cfg.Brokers) > 0 && cfg.Brokers[0] != "" {
		broker = cfg.Brokers[0]
	}
	return messaging.KafkaConfig{
		Brokers:              []string{broker},
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
	repo, err := initStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	messagingConfig := getKafkaSettings(cfg.Messaging)
	kafkaClient, err := messaging.NewKafkaClient(messagingConfig)
	if err != nil {
		return nil, err
	}
	cloudinarySvc, err := infrastructure.NewCloudinaryService(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		return nil, fmt.Errorf("init cloudinary service: %w", err)
	}
	httpRouter := setupHttpRouter(cfg, repo, kafkaClient, cloudinarySvc)

	return &container{
		repo:          repo,
		server:        server.New(getServerConfig(cfg), httpRouter),
		messaging:     kafkaClient,
		cloudinarySvc: cloudinarySvc,
	}, nil
}
func initStorage(cfg config.Config) (domain.SellerRepository, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres init error: %w", err)
	}

	return repo, nil
}

func setupHttpRouter(cfg config.Config, r domain.SellerRepository, m domain.Messaging, c domain.ImageService) server.RouteRegistrar {

	httpHandlers := httptransport.NewHandlers(r, m, c)
	return httptransport.NewRouter(httpHandlers)
}

func getServerConfig(cfg config.Config) server.Config {
	return server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
