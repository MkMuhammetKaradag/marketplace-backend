package messaging

import (
	"time"
)

type KafkaConfig struct {
	Brokers              []string
	GroupID              string
	Topic                string
	ClientID             string
	QueueDurable         bool
	QueueAutoDelete      bool
	EnableRetry          bool
	MaxRetries           int
	RetryTopic           string
	DLQTopic             string
	ConnectionTimeout    time.Duration
	ServiceType          ServiceType
	AllowedMessageTypes  map[ServiceType][]MessageType
	CriticalMessageTypes []MessageType
}

func NewDefaultConfig(kafkaBrokers []string) KafkaConfig {
	if kafkaBrokers == nil || len(kafkaBrokers) == 0 {
		kafkaBrokers = []string{"localhost:9092"}
	}

	return KafkaConfig{
		Brokers:              kafkaBrokers,
		Topic:                "main-events",
		RetryTopic:           "main-events-retry",
		DLQTopic:             "main-events-dlq",
		ServiceType:          ServiceType_UNKNOWN_SERVICE,
		EnableRetry:          true,
		MaxRetries:           3,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []MessageType{},
	}
}
