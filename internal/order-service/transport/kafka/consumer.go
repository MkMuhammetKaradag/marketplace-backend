package kafka

import (
	"context"
	"fmt"
	"log"
	"marketplace/internal/order-service/config"
	"marketplace/internal/order-service/domain"
	"marketplace/pkg/messaging"
	eventsProto "marketplace/pkg/proto/events"
	"time"
)

type Consumer struct {
	client   *messaging.KafkaClient
	handlers map[eventsProto.MessageType]domain.MessageHandler
	cfg      messaging.KafkaConfig
}

func NewConsumer(cfg config.MessagingConfig, handlers map[eventsProto.MessageType]domain.MessageHandler) (*Consumer, error) {
	kafkaConfig := createKafkaConfig(cfg)
	client, err := messaging.NewKafkaClient(kafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("init kafka messaging: %w", err)
	}

	return &Consumer{
		client:   client,
		handlers: handlers,
		cfg:      kafkaConfig,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	messageRouter := func(ctx context.Context, msg *eventsProto.Message) error {
		handler, ok := c.handlers[msg.Type]
		if !ok {
			return nil
		}
		return handler.Handle(ctx, msg)
	}
	ctx, cancel := context.WithCancel(ctx)
	// Ana Consumer
	go func() {
		log.Println("Starting Kafka consumer for main-events...")
		groupID := c.cfg.ServiceType.String() + "-main-group"
		topic := c.cfg.Topic
		if err := c.client.ConsumeMessages(ctx, messageRouter, &topic, &groupID); err != nil {
			log.Printf("Main consumer error: %v", err)
			cancel()
		}
	}()

	go func() {
		log.Println("Starting Kafka consumer for retry-events...")
		// Aynı KafkaClient'ı kullanarak, sadece farklı bir topic ve grup adı veriyoruz.
		groupID := c.cfg.ServiceType.String() + "-retry-group"
		topic := c.cfg.RetryTopic
		if err := c.client.ConsumeMessages(ctx, messageRouter, &topic, &groupID); err != nil {
			log.Printf("Retry consumer error: %v", err)
			cancel()
		}
	}()

	// DLQ Consumer
	go func() {
		log.Println("Starting DLQ recovery consumer...")
		if err := c.client.ConsumeDLQWithRecovery(ctx, messageRouter); err != nil {
			log.Printf("DLQ consumer error: %v", err)
			cancel()
		}
	}()
}

func (c *Consumer) Client() domain.Messaging {
	return c.client
}

func createKafkaConfig(cfg config.MessagingConfig) messaging.KafkaConfig {
	broker := cfg.Brokers[0]
	if broker == "" {
		broker = "localhost:9092"
	}
	kafkaBrokers := []string{broker}
	return messaging.KafkaConfig{
		Brokers:               kafkaBrokers,
		Topic:                 "main-events",
		RetryTopic:            "retry-events",
		DLQTopic:              "dlq-events",
		ServiceType:           eventsProto.ServiceType_ORDER_SERVICE,
		EnableRetry:           true,
		MaxRetries:            10,
		ConnectionTimeout:     10 * time.Second,
		MaxConcurrentHandlers: 10,
		AllowedMessageTypes: map[eventsProto.ServiceType][]eventsProto.MessageType{
			eventsProto.ServiceType_ORDER_SERVICE: {eventsProto.MessageType_PAYMENT_SUCCESSFUL},
		},
		CriticalMessageTypes: []eventsProto.MessageType{eventsProto.MessageType_PAYMENT_SUCCESSFUL},
	}
}
