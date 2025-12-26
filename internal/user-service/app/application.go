// internal/user-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/user-service/config"
	"marketplace/internal/user-service/domain"
	"marketplace/internal/user-service/infrastructure"
	"marketplace/internal/user-service/pkg/graceful"
	"marketplace/internal/user-service/repository/postgres"
	"marketplace/internal/user-service/repository/session"
	"marketplace/internal/user-service/server"
	grpctransport "marketplace/internal/user-service/transport/grpc"
	httptransport "marketplace/internal/user-service/transport/http"
	"marketplace/internal/user-service/transport/kafka"
	messaginghandler "marketplace/internal/user-service/transport/messaging"

	"time"
)

type App struct {
	cfg           config.Config
	server        *server.Server
	repository    domain.UserRepository
	sessionRepo   domain.SessionRepository
	messaging     domain.Messaging
	consumer      *kafka.Consumer
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
		sessionRepo:   container.sessionRepo,
		messaging:     container.messaging,
		consumer:      container.consumer,
		cloudinarySvc: container.cloudinarySvc,
	}, nil
}

func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer
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
	repo          domain.UserRepository
	sessionRepo   domain.SessionRepository
	server        *server.Server
	messaging     domain.Messaging
	consumer      *kafka.Consumer
	cloudinarySvc domain.ImageService
}

func buildContainer(cfg config.Config) (*container, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres repository: %w", err)
	}
	sessionRepo, err := session.NewSessionRepository(cfg)
	if err != nil {
		return nil, fmt.Errorf("init redis session manager: %w", err)
	}

	messsagingHnadlers := messaginghandler.SetupMessageHandlers(repo)

	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, messsagingHnadlers)
	if err != nil {
		return nil, fmt.Errorf("init kafka consumer: %w", err)
	}

	cloudinarySvc, err := infrastructure.NewCloudinaryService(cfg.Cloudinary.CloudName, cfg.Cloudinary.APIKey, cfg.Cloudinary.APISecret)
	if err != nil {
		return nil, fmt.Errorf("init cloudinary service: %w", err)
	}
	userService := domain.NewUserService(repo)
	httpHandlers := httptransport.NewHandlers(userService, repo, sessionRepo, cloudinarySvc)

	router := httptransport.NewRouter(httpHandlers)

	serverCfg := server.Config{
		Port:         cfg.Server.Port,
		IdleTimeout:  5 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		GrpcPort:     cfg.Server.GrpcPort,
	}

	grpcHandler := grpctransport.NewAuthGrpcHandler(sessionRepo)
	httpServer := server.New(serverCfg, router, grpcHandler)

	return &container{
		repo:          repo,
		server:        httpServer,
		sessionRepo:   sessionRepo,
		messaging:     kafkaConsumer.Client(),
		consumer:      kafkaConsumer,
		cloudinarySvc: cloudinarySvc,
	}, nil
}
