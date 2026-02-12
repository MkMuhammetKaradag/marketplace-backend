// internal/payment-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/payment-service/config"
	"marketplace/internal/payment-service/domain"
	"marketplace/internal/payment-service/infrastructure/payment"
	"marketplace/internal/payment-service/pkg/graceful"
	"marketplace/internal/payment-service/repository/postgres"
	"marketplace/internal/payment-service/server"
	grpctransport "marketplace/internal/payment-service/transport/grpc"
	httptransport "marketplace/internal/payment-service/transport/http"
	"marketplace/internal/payment-service/transport/kafka"
	messaginghandler "marketplace/internal/payment-service/transport/messaging"
	"marketplace/pkg/messaging"
	eventsProto "marketplace/pkg/proto/events"
	"time"
)

type App struct {
	cfg      config.Config
	server   *server.Server
	consumer *kafka.Consumer
	repo     domain.PaymentRepository // ✅ Repo eklenmeli (Close için)
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
		repo:     container.repo,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go a.consumer.Start(ctx)

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)

	log.Printf("starting payment-service on %s (gRPC: %s)", a.cfg.Server.Port, a.cfg.Server.GrpcPort)
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.repo.Close()
}

type container struct {
	server    *server.Server
	consumer  *kafka.Consumer
	messaging domain.Messaging
	repo      domain.PaymentRepository
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
	// 1. Storage
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	// 2. External Services (Stripe vb.)
	stripeSvc := payment.NewStripeService(cfg.Stripe.SecretKey, cfg.Stripe.WebhookSecret)

	// 3. Messaging (Kafka Client & Consumer)
	mCfg := createMessagingConfig(cfg.Messaging)
	kafkaClient, err := messaging.NewKafkaClient(mCfg)
	if err != nil {
		return nil, err
	}

	msgHandlers := messaginghandler.SetupMessageHandlers() // ✅ Repo ve Svc ekle
	consumer, err := kafka.NewConsumer(cfg.Messaging, msgHandlers)
	if err != nil {
		return nil, err
	}

	// 4. Transport (HTTP & gRPC)

	h := httptransport.NewHandlers(repo, stripeSvc, kafkaClient)
	router := httptransport.NewRouter(h)
	grpcHandler := grpctransport.NewPaymentGrpcHandler(repo, stripeSvc)

	return &container{
		repo:      repo,
		server:    server.New(getServerConfig(cfg), router, grpcHandler),
		consumer:  consumer,
		messaging: kafkaClient,
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
