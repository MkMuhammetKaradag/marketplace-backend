// internal/product-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/product-service/config"
	"marketplace/internal/product-service/domain"
	"marketplace/internal/product-service/infrastructure/ai"
	"marketplace/internal/product-service/infrastructure/img"
	"marketplace/internal/product-service/infrastructure/worker"
	"marketplace/internal/product-service/repository/postgres"
	"marketplace/internal/product-service/server"
	grpctransport "marketplace/internal/product-service/transport/grpc"
	httptransport "marketplace/internal/product-service/transport/http"
	"marketplace/internal/product-service/transport/kafka"
	messaginghandler "marketplace/internal/product-service/transport/messaging"
	"marketplace/pkg/messaging"

	"time"

	pb "marketplace/pkg/proto/events"

	"github.com/hibiken/asynq"
)

type App struct {
	cfg        config.Config
	server     *server.Server
	repository domain.ProductRepository
	//messaging     domain.Messaging
	consumer *kafka.Consumer
	//cloudinarySvc domain.ImageService
	//aiProvider    domain.AiProvider
	asynqClient *asynq.Client
	// worker        domain.Worker
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
		//messaging:   container.messaging,
		consumer:    container.consumer,
		asynqClient: container.asynqClient,
		// worker:      container.worker,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Kafka Consumer'ı başlat
	a.consumer.Start(ctx)

	// 2. Transactional Outbox Relay'i başlat (Eğer Product DB'de outbox varsa)
	// go a.startOutboxRelay(ctx)

	log.Printf("starting product-service on %s (gRPC: %s)", a.cfg.Server.Port, a.cfg.Server.GrpcPort)

	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return nil
}

type container struct {
	repo          domain.ProductRepository
	server        *server.Server
	messaging     domain.Messaging
	consumer      *kafka.Consumer
	cloudinarySvc domain.ImageService
	aiProvider    domain.AiProvider
	asynqClient   *asynq.Client
	worker        domain.Worker
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
		ServiceType:          pb.ServiceType_PRODUCT_SERVICE,
		EnableRetry:          true,
		MaxRetries:           10,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []pb.MessageType{pb.MessageType_PRODUCT_PRICE_UPDATED},
	}
}
func buildContainer(cfg config.Config) (*container, error) {
	// Veritabanı
	repo, err := initStorage(cfg)
	if err != nil {
		return nil, err
	}

	// Dış Servisler (AI, Cloudinary)
	aiProvider := ai.NewOllamaProvider()
	cloudinarySvc, err := img.NewCloudinaryService(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		return nil, err
	}

	// Redis & Asynq Yapılandırması
	redisOpt := asynq.RedisClientOpt{Addr: "localhost:6379", DB: 2}
	asynqClient := asynq.NewClient(redisOpt)
	wrk := worker.NewWorker(asynqClient)

	// Buradaki Processor'ı ayrı bir goroutine'de başlatmak yerine container'a ekleyebiliriz
	// veya Start metodunda tetikleyebiliriz.
	processor := worker.NewTaskProcessor(redisOpt, repo, cloudinarySvc)
	go func() {
		if err := processor.Start(); err != nil {
			log.Printf("Task Processor error: %v", err)
		}
	}()

	// Messaging (Kafka)
	msgHandlers := messaginghandler.SetupMessageHandlers(repo)
	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, msgHandlers)
	if err != nil {
		return nil, err
	}

	// Transport (HTTP & gRPC)
	productService := domain.NewProductService(repo)
	httpRouter := setupHttpRouter(cfg, productService, repo, cloudinarySvc, aiProvider, wrk, kafkaConsumer.Client())
	grpcHandler := grpctransport.NewProductGrpcHandler(repo)

	return &container{
		repo:        repo,
		server:      server.New(getServerConfig(cfg), httpRouter, grpcHandler),
		messaging:   kafkaConsumer.Client(),
		consumer:    kafkaConsumer,
		asynqClient: asynqClient,
		aiProvider:  aiProvider,
		worker:      wrk,
	}, nil
}
func initStorage(cfg config.Config) (domain.ProductRepository, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgres init error: %w", err)
	}

	return repo, nil
}
func setupHttpRouter(cfg config.Config, p domain.ProductService, r domain.ProductRepository, c domain.ImageService, a domain.AiProvider, w domain.Worker, m domain.Messaging) server.RouteRegistrar {

	httpHandlers := httptransport.NewHandlers(p, r, c, a, w, m)
	return httptransport.NewRouter(httpHandlers)
}

func getServerConfig(cfg config.Config) server.Config {
	return server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		GrpcPort:     cfg.Server.GrpcPort,
	}
}
