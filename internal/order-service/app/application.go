// internal/order-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/order-service/config"
	"marketplace/internal/order-service/domain"
	"marketplace/internal/order-service/grpc_client"
	"marketplace/internal/order-service/pkg/graceful"
	"marketplace/internal/order-service/repository/postgres"
	"marketplace/internal/order-service/server"
	httptransport "marketplace/internal/order-service/transport/http"
	"marketplace/internal/order-service/transport/kafka"
	messaginghandler "marketplace/internal/order-service/transport/messaging"
	"marketplace/pkg/messaging"
	eventsProto "marketplace/pkg/proto/events"
	"time"
)

type App struct {
	cfg        config.Config
	server     *server.Server
	repository domain.OrderRepository
	consumer   *kafka.Consumer
}

func NewApp(cfg config.Config) (*App, error) {
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("bootstrap failed: %w", err)
	}

	return &App{
		cfg:        cfg,
		server:     container.server,
		repository: container.repository,
		consumer:   container.consumer,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)
	go a.consumer.Start(ctx)
	log.Printf("starting user-service on %s", a.server.Address())

	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.repository.Close()
}

type container struct {
	server     *server.Server
	repository domain.OrderRepository
	consumer   *kafka.Consumer
	messaging  domain.Messaging
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
		ServiceType:          eventsProto.ServiceType_ORDER_SERVICE,
		EnableRetry:          true,
		MaxRetries:           10,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []eventsProto.MessageType{eventsProto.MessageType_PAYMENT_SUCCESSFUL},
	}
}
func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	grpcProductAddress := fmt.Sprintf("localhost:%s", cfg.Server.GrpcProductPort)
	grpcProductClient, err := grpc_client.NewProductClient(grpcProductAddress)
	if err != nil {
		log.Fatalf("failed to initialise gRPC product client: %v", err)
	}

	grpcBasketAddress := fmt.Sprintf("localhost:%s", cfg.Server.GrpcBasketPort)
	grpcBasketClient, err := grpc_client.NewBasketClient(grpcBasketAddress)
	if err != nil {
		log.Fatalf("failed to initialise gRPC basket client: %v", err)
	}

	grpcPaymentAddress := fmt.Sprintf("localhost:%s", cfg.Server.GrpcPaymentPort)
	grpcPaymentClient, err := grpc_client.NewPaymentClient(grpcPaymentAddress)
	if err != nil {
		log.Fatalf("failed to initialise gRPC basket client: %v", err)
	}

	messagingConfig := createMessagingConfig(cfg.Messaging)
	messaging, err := messaging.NewKafkaClient(messagingConfig)
	if err != nil {
		return nil, fmt.Errorf("init kafka messaging: %w", err)
	}
	messsagingHnadlers := messaginghandler.SetupMessageHandlers(repo)
	httpHandlers := httptransport.NewHandlers(repo, grpcProductClient, grpcBasketClient, grpcPaymentClient, messaging)
	router := httptransport.NewRouter(httpHandlers)
	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, messsagingHnadlers)
	if err != nil {
		return nil, fmt.Errorf("init kafka consumer: %w", err)
	}

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	httpServer := server.New(serverCfg, router, grpcBasketClient, grpcProductClient)

	return &container{

		server:     httpServer,
		repository: repo,
		consumer:   kafkaConsumer,
		messaging:  kafkaConsumer.Client(),
	}, nil
}
