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
	"marketplace/pkg/messaging"
	pb "marketplace/pkg/proto/events"

	"time"

	"google.golang.org/protobuf/proto"
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
	go a.startOutboxRelay(ctx)
	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)

	log.Printf("starting user-service on %s", a.server.Address())
	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	log.Println("server stopped, closing repository")
	return a.repository.Close()
}

func (a *App) startOutboxRelay(ctx context.Context) {
	// Her 2 saniyede bir kontrol et (İhtiyaca göre 500ms de olabilir)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	log.Println("Outbox Relay worker started...")

	for {
		select {
		case <-ctx.Done(): // Uygulama kapanırken işçiyi de durdururuz
			log.Println("Outbox Relay worker stopping...")
			return
		case <-ticker.C:
			a.processOutboxMessages(ctx)
		}
	}
}
func (a *App) processOutboxMessages(ctx context.Context) {
	// 1. Bekleyen mesajları çek (Örn: her seferinde 10 tane)
	messages, err := a.repository.GetPendingOutboxMessages(ctx, 10)
	if err != nil {
		log.Printf("Relay: error fetching messages: %v", err)
		return
	}

	if len(messages) == 0 {
		return // Bekleyen mesaj yoksa bir sonraki ticker'ı bekle
	}

	for _, msg := range messages {
		// 2. Mesajı Kafka'ya gönder
		// msg.Payload zaten []byte olduğu için doğrudan gönderebiliriz
		// Eğer PublishMessage metodun pb.Message istiyorsa Unmarshal etmen gerekir

		var protoMsg pb.Message
		if err := proto.Unmarshal(msg.Payload, &protoMsg); err != nil {
			log.Printf("Relay: unmarshal error: %v", err)
			continue
		}

		err := a.messaging.PublishMessage(ctx, &protoMsg)
		if err != nil {
			log.Printf("Relay: failed to publish message %s: %v", msg.ID, err)
			continue // Bir sonrakine geç, bu PENDING kalmaya devam edecek
		}
		// 3. Başarılıysa DB'de PROCESSED yap
		if err := a.repository.MarkOutboxAsProcessed(ctx, msg.ID); err != nil {
			log.Printf("Relay: failed to mark message %s as processed: %v", msg.ID, err)
		}
	}
}

type container struct {
	repo          domain.UserRepository
	sessionRepo   domain.SessionRepository
	server        *server.Server
	messaging     domain.Messaging
	consumer      *kafka.Consumer
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
		ServiceType:          pb.ServiceType_USER_SERVICE,
		EnableRetry:          true,
		MaxRetries:           3,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []pb.MessageType{pb.MessageType_USER_CREATED},
	}
}
func buildContainer(cfg config.Config) (*container, error) {

	repo, sessionRepo, err := initStorage(cfg)
	if err != nil {
		return nil, err
	}

	cloudinarySvc, err := infrastructure.NewCloudinaryService(
		cfg.Cloudinary.CloudName,
		cfg.Cloudinary.APIKey,
		cfg.Cloudinary.APISecret,
	)
	if err != nil {
		return nil, fmt.Errorf("cloudinary init failed: %w", err)
	}

	messagingHandlers := messaginghandler.SetupMessageHandlers(repo)
	kafkaConsumer, err := kafka.NewConsumer(cfg.Messaging, messagingHandlers)
	if err != nil {
		return nil, fmt.Errorf("kafka init failed: %w", err)
	}
	msgClient := kafkaConsumer.Client()

	httpRouter := setupHttpRouter(cfg, repo, sessionRepo, msgClient, cloudinarySvc)
	grpcHandler := grpctransport.NewAuthGrpcHandler(sessionRepo)

	return &container{
		repo:          repo,
		sessionRepo:   sessionRepo,
		consumer:      kafkaConsumer,
		cloudinarySvc: cloudinarySvc,
		server:        server.New(getServerConfig(cfg), httpRouter, grpcHandler),
		messaging:     msgClient,
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

func initStorage(cfg config.Config) (domain.UserRepository, domain.SessionRepository, error) {
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("postgres init error: %w", err)
	}

	sessionRepo, err := session.NewSessionRepository(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("redis init error: %w", err)
	}

	return repo, sessionRepo, nil
}

func setupHttpRouter(cfg config.Config, r domain.UserRepository, s domain.SessionRepository, m domain.Messaging, c domain.ImageService) server.RouteRegistrar {
	userService := domain.NewUserService(r)

	httpHandlers := httptransport.NewHandlers(userService, r, s, m, c)
	return httptransport.NewRouter(httpHandlers)
}
