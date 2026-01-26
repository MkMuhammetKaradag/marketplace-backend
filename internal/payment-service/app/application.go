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
	"marketplace/internal/payment-service/transport/kafka"
	messaginghandler "marketplace/internal/payment-service/transport/messaging"
	"marketplace/internal/user-service/domain"
	"marketplace/pkg/messaging"
	eventsProto "marketplace/pkg/proto/events"
	"time"
)

type App struct {
	cfg      config.Config
	server   *server.Server
	consumer *kafka.Consumer
}

func NewApp(cfg config.Config) (*App, error) {
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return &App{
		cfg:      cfg,
		server:   container.server,
		consumer: container.consumer,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)
	go a.consumer.Start(ctx)

	log.Printf("starting payment-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return nil
}

type container struct {
	server    *server.Server
	consumer  *kafka.Consumer
	messaging domain.Messaging
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
		ServiceType:          eventsProto.ServiceType_PAYMENT_SERVICE,
		EnableRetry:          true,
		MaxRetries:           10,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []eventsProto.MessageType{eventsProto.MessageType_ORDER_CREATED},
	}
}
func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}
	stripeService := payment.NewStripeService(cfg.Stripe.SecretKey, cfg.Stripe.WebhookSecret)

	messagingConfig := createMessagingConfig(cfg.Messaging)
	messaging, err := messaging.NewKafkaClient(messagingConfig)
	if err != nil {
		return nil, fmt.Errorf("init kafka messaging: %w", err)
	}
	messsagingHnadlers := messaginghandler.SetupMessageHandlers()
	httpHandlers := httptransport.NewHandlers(stripeService, messaging)
	router := httptransport.NewRouter(httpHandlers)
	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, messsagingHnadlers)
	if err != nil {
		return nil, fmt.Errorf("init kafka consumer: %w", err)
	}
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

		server:    httpServer,
		consumer:  kafkaConsumer,
		messaging: kafkaConsumer.Client(),
	}, nil
}
