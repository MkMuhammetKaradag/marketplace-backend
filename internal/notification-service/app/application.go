// internal/notification-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/notification-service/config"
	"marketplace/internal/notification-service/domain"
	email "marketplace/internal/notification-service/infrastructure"
	"marketplace/internal/notification-service/pkg/graceful"
	"marketplace/internal/notification-service/repository/postgres"
	"marketplace/internal/notification-service/server"
	httptransport "marketplace/internal/notification-service/transport/http"
	"marketplace/internal/notification-service/transport/kafka"
	messaginghandler "marketplace/internal/notification-service/transport/messaging"
	"marketplace/pkg/messaging"
	pb "marketplace/pkg/proto/events"
	"time"
)

type App struct {
	cfg        config.Config
	server     *server.Server
	consumer   *kafka.Consumer
	repository domain.NotificationRepository
}

func NewApp(cfg config.Config) (*App, error) {
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return &App{
		cfg:        cfg,
		consumer:   container.consumer,
		server:     container.server,
		repository: container.repo,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.consumer.Start(ctx)
	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)
	log.Printf("starting user-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.repository.Close()
}

type container struct {
	repo     domain.NotificationRepository
	server   *server.Server
	consumer *kafka.Consumer
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
		ServiceType:          pb.ServiceType_NOTIFICATION_SERVICE,
		EnableRetry:          true,
		MaxRetries:           10,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []pb.MessageType{pb.MessageType_USER_ACTIVATION_EMAIL},
	}
}
func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}
	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	resendProvider := email.NewResendProvider(cfg.Email.ApiKey)
	messsagingHnadlers := messaginghandler.SetupMessageHandlers(resendProvider, repo)

	messagingConfig := createMessagingConfig(cfg.Messaging)
	messaging, err := messaging.NewKafkaClient(messagingConfig)
	if err != nil {
		return nil, fmt.Errorf("init kafka messaging: %w", err)
	}

	httpHandlers := httptransport.NewHandlers(messaging)
	router := httptransport.NewRouter(httpHandlers)

	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, messsagingHnadlers)
	if err != nil {
		return nil, fmt.Errorf("init kafka consumer: %w", err)
	}
	httpServer := server.New(serverCfg, router)

	return &container{
		repo:     repo,
		consumer: kafkaConsumer,
		server:   httpServer,
	}, nil
}
