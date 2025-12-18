package messaging

import (
	pb "marketplace/pkg/proto/events"
	"time"
)

type KafkaConfig struct {
	Brokers               []string
	GroupID               string
	Topic                 string
	ClientID              string
	QueueDurable          bool
	QueueAutoDelete       bool
	EnableRetry           bool
	MaxRetries            int
	MaxConcurrentHandlers int
	RetryTopic            string
	DLQTopic              string
	ConnectionTimeout     time.Duration
	ServiceType           pb.ServiceType

	AllowedMessageTypes  map[pb.ServiceType][]pb.MessageType
	CriticalMessageTypes []pb.MessageType
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
		ServiceType:          pb.ServiceType_UNKNOWN_SERVICE,
		EnableRetry:          true,
		MaxRetries:           3,
		ConnectionTimeout:    10 * time.Second,
		CriticalMessageTypes: []pb.MessageType{},
	}
}
