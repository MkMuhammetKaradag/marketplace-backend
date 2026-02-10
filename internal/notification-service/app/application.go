// internal/notification-service/app/application.go
package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/notification-service/config"
	"marketplace/internal/notification-service/domain"
	email "marketplace/internal/notification-service/infrastructure/email"
	template_manager "marketplace/internal/notification-service/infrastructure/template_manager"
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

type container struct {
	repo     domain.NotificationRepository
	server   *server.Server
	consumer *kafka.Consumer
}
type App struct {
	cfg        config.Config
	server     *server.Server
	consumer   *kafka.Consumer
	repository domain.NotificationRepository
}

func NewApp(cfg config.Config) (*App, error) {
	// buildContainer'ı çağırıp tüm bağımlılıkları (repo, kafka, server) hazırlıyoruz
	container, err := buildContainer(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build container: %w", err)
	}

	// App struct'ını doldurup geri dönüyoruz
	return &App{
		cfg:        cfg,
		server:     container.server,
		consumer:   container.consumer,
		repository: container.repo,
	}, nil
}
func (a *App) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.consumer.Start(ctx)

	// Log mesajını dinamikleştirin (user-service değil notification-service)
	log.Printf("Starting Notification Service on %s", a.server.Address())

	go graceful.WaitForShutdown(a.server.FiberApp(), 5*time.Second, ctx)

	if err := a.server.Start(); err != nil {
		return fmt.Errorf("server exited with error: %w", err)
	}

	return a.repository.Close()
}

func buildContainer(cfg config.Config) (*container, error) {
	// 1. Veritabanı Başlatma
	repo, err := postgres.NewRepository(cfg)
	if err != nil {
		return nil, err
	}

	// 2. Altyapı Servisleri (Infrastructure)
	// EmailProvider interface üzerinden tanımlanmalı
	emailProvider := email.NewResendProvider(cfg.Email.ApiKey)


	templateMgr := template_manager.NewTemplateManager("templates")

	// 3. İş Mantığı ve Handlerlar (Transport/Messaging)
	// Bu kısım çok şişiyorsa bir 'Dependency Registry' oluşturulabilir
	handlers := messaginghandler.SetupMessageHandlers(emailProvider, templateMgr, repo)

	// 4. Messaging (Kafka) Başlatma
	messagingConfig := getKafkaSettings(cfg.Messaging)
	kafkaClient, err := messaging.NewKafkaClient(messagingConfig)
	if err != nil {
		return nil, err
	}

	// 5. Taşıma Katmanları (HTTP & Kafka Consumer)
	httpRouter := setupRouter(kafkaClient)

	consumer, err := kafka.NewConsumer(cfg.Messaging, handlers)
	if err != nil {
		return nil, err
	}

	return &container{
		repo:     repo,
		consumer: consumer,
		server:   server.New(getServerConfig(cfg), httpRouter),
	}, nil
}

// YARDIMCI FONKSİYONLAR - Kodun kalabalığını aşağıya taşıyoruz
func getKafkaSettings(cfg config.MessagingConfig) messaging.KafkaConfig {
	broker := "localhost:29092"
	if len(cfg.Brokers) > 0 && cfg.Brokers[0] != "" {
		broker = cfg.Brokers[0]
	}

	return messaging.KafkaConfig{
		Brokers:              []string{broker},
		Topic:                "main-events",
		RetryTopic:           "main-events-retry",
		DLQTopic:             "main-events-dlq",
		ServiceType:          pb.ServiceType_NOTIFICATION_SERVICE,
		EnableRetry:          true,
		MaxRetries:           10,
		CriticalMessageTypes: []pb.MessageType{pb.MessageType_USER_ACTIVATION_EMAIL},
	}
}

func setupRouter(msgClient domain.Messaging) server.RouteRegistrar {
	httpHandlers := httptransport.NewHandlers(msgClient)
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
