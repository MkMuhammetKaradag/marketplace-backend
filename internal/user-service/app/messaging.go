package app

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/user-service/config"
	"marketplace/internal/user-service/domain"
	"marketplace/pkg/messaging"
	"time"
)

func createMessagingConfig(cfg config.MessagingConfig) messaging.KafkaConfig {
	broker := cfg.Brokers[0]
	if broker == "" {
		broker = "localhost:9092"
	}
	kafkaBrokers := []string{broker}
	return messaging.KafkaConfig{
		Brokers:           kafkaBrokers,
		Topic:             "main-events",
		RetryTopic:        "retry-events",
		DLQTopic:          "dlq-events",
		ServiceType:       messaging.ServiceType_USER_SERVICE,
		EnableRetry:       true,
		MaxRetries:        3,
		ConnectionTimeout: 10 * time.Second,
		AllowedMessageTypes: map[messaging.ServiceType][]messaging.MessageType{
			messaging.ServiceType_USER_SERVICE: {
				messaging.MessageType_SELLER_APPROVED,
				messaging.MessageType_SELLER_REJECTED,
			},
		},
		CriticalMessageTypes: []messaging.MessageType{messaging.MessageType_USER_CREATED},
	}
}

func SetupMessaging(handlers map[messaging.MessageType]domain.MessageHandler, cfg config.Config) (domain.Messaging, error) {

	messageRouter := func(ctx context.Context, msg *messaging.Message) error {
		handler, ok := handlers[msg.Type]
		fmt.Println("messaj geldi-msg.type", msg.Type)
		if !ok {
			return nil
		}
		return handler.Handle(ctx, msg)
	}
	config := createMessagingConfig(cfg.Messaging)

	messaging, err := messaging.NewKafkaClient(config)
	if err != nil {
		return nil, fmt.Errorf("init kafka messaging: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	// Consumer'ları başlat
	// Ana Consumer
	go func() {
		log.Println("Starting Kafka consumer for main-events...")
		groupID := config.ServiceType.String() + "-main-group"
		topic := config.Topic
		if err := messaging.ConsumeMessages(ctx, messageRouter, &topic, &groupID); err != nil {
			log.Printf("Main consumer error: %v", err)
			cancel()
		}
	}()

	// // Retry Consumer
	// go func() {
	// 	log.Println("Starting Kafka consumer for retry-events...")
	// 	// Aynı KafkaClient'ı kullanarak, sadece farklı bir topic ve grup adı veriyoruz.
	// 	groupID := config.ServiceType.String() + "-retry-group"
	// 	topic := config.RetryTopic
	// 	if err := messaging.ConsumeMessages(ctx, messageRouter, &topic, &groupID); err != nil {
	// 		log.Printf("Retry consumer error: %v", err)
	// 		cancel()
	// 	}
	// }()

	// // DLQ Consumer
	// go func() {
	// 	log.Println("Starting DLQ recovery consumer...")
	// 	if err := messaging.ConsumeDLQWithRecovery(ctx, messageRouter); err != nil {
	// 		log.Printf("DLQ consumer error: %v", err)
	// 	}
	// }()

	return messaging, nil
}
