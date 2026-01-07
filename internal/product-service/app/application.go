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
	httptransport "marketplace/internal/product-service/transport/http"
	"marketplace/internal/product-service/transport/kafka"
	messaginghandler "marketplace/internal/product-service/transport/messaging"

	"time"

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
	defer a.asynqClient.Close()
	a.consumer.Start(ctx)
	log.Printf("starting user-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.repository.Close()
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

func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}

	aiProvider := ai.NewOllamaProvider()
	cloudinarySvc, err := img.NewCloudinaryService(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		return nil, fmt.Errorf("init cloudinary service: %w", err)
	}
	productService := domain.NewProductService(repo)

	redisOpt := asynq.RedisClientOpt{
		Addr:     "localhost:6379",
		Password: "",
		DB:       2,
	}

	// 2. Asynq Client'ı oluştur (Bu senin 'want *asynq.Client' kısmın)
	asynqClient := asynq.NewClient(redisOpt)

	wrk := worker.NewWorker(asynqClient)
	httpHandlers := httptransport.NewHandlers(productService, repo, cloudinarySvc, aiProvider, wrk)
	messsagingHnadlers := messaginghandler.SetupMessageHandlers(repo)
	router := httptransport.NewRouter(httpHandlers)

	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, messsagingHnadlers)
	if err != nil {
		return nil, fmt.Errorf("init kafka consumer: %w", err)
	}
	processor := worker.NewTaskProcessor(redisOpt, repo, cloudinarySvc)

	go func() {
		if err := processor.Start(); err != nil {
			log.Fatalf("Worker başlatılamadı: %v", err)
		}
	}()
	//defer asynqClient.Close()
	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		GrpcPort:     cfg.Server.GrpcPort,
	}

	httpServer := server.New(serverCfg, router)

	return &container{
		repo:        repo,
		server:      httpServer,
		messaging:   kafkaConsumer.Client(),
		consumer:    kafkaConsumer,
		asynqClient: asynqClient,
	}, nil
}
