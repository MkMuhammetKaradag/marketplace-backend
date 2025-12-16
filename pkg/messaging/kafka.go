package messaging

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaClient struct {
	config      KafkaConfig
	producer    *kafka.Writer
	mu          sync.Mutex
	closed      bool
	ServiceType ServiceType
}

func NewKafkaClient(config KafkaConfig) (*KafkaClient, error) {
	kc := &KafkaClient{
		config:      config,
		ServiceType: config.ServiceType,
	}

	if err := kc.createTopicsIfNotExists(); err != nil {
		log.Println("Failed to create topics:", err)
	}

	producer := &kafka.Writer{
		Addr:         kafka.TCP(config.Brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		RequiredAcks: kafka.RequireAll,
		BatchSize:    100,
		BatchBytes:   1048576,
		BatchTimeout: 1 * time.Second,
		MaxAttempts:  3,

		// Logger:      kafka.LoggerFunc(quietLogger.Printf),
		// ErrorLogger: kafka.LoggerFunc(log.Printf),
	}
	kc.producer = producer

	log.Printf("Kafka Client initialized for service: %s, main topic: %s", config.ServiceType.String(), config.Topic)
	return kc, nil

}
func (kc *KafkaClient) createTopicsIfNotExists() error {
	conn, err := kafka.Dial("tcp", kc.config.Brokers[0])
	if err != nil {
		return fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("failed to get controller: %w", err)
	}

	controllerConn, err := kafka.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return fmt.Errorf("failed to connect to controller: %w", err)
	}
	defer controllerConn.Close()

	// Oluşturulacak topic'ler
	topics := []kafka.TopicConfig{
		{
			Topic:             kc.config.Topic, // main-events
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	// Retry topic varsa ekle
	if kc.config.EnableRetry && kc.config.RetryTopic != "" {
		topics = append(topics, kafka.TopicConfig{
			Topic:             kc.config.RetryTopic,
			NumPartitions:     3,
			ReplicationFactor: 1,
		})
	}

	// DLQ topic varsa ekle
	if kc.config.DLQTopic != "" {
		topics = append(topics, kafka.TopicConfig{
			Topic:             kc.config.DLQTopic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	err = controllerConn.CreateTopics(topics...)
	if err != nil {
		// Topic zaten varsa hata vermemeli
		if !strings.Contains(err.Error(), "already exists") {
			return fmt.Errorf("failed to create topics: %w", err)
		}
		log.Println("Some topics already exist, continuing...")
	}

	// Oluşturulan topic'leri logla
	for _, topic := range topics {
		log.Printf("Topic '%s' created or already exists", topic.Topic)
	}

	return nil
}
func (kc *KafkaClient) Close() error {
	kc.mu.Lock()
	defer kc.mu.Unlock()

	kc.closed = true

	if kc.producer != nil {
		if err := kc.producer.Close(); err != nil {
			return fmt.Errorf("failed to close Kafka producer: %w", err)
		}
	}
	log.Println("Kafka Client closed.")
	return nil
}
